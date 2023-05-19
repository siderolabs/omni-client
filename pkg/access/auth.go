// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package access

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	pgpcrypto "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/pkg/browser"
	authcli "github.com/siderolabs/go-api-signature/pkg/client/auth"
	"github.com/siderolabs/go-api-signature/pkg/client/interceptor"
	"github.com/siderolabs/go-api-signature/pkg/message"
	"github.com/siderolabs/go-api-signature/pkg/pgp"
	"github.com/siderolabs/go-api-signature/pkg/pgp/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/siderolabs/omni-client/pkg/client/omni"
	authres "github.com/siderolabs/omni-client/pkg/omni/resources/auth"
	"github.com/siderolabs/omni-client/pkg/version"
)

// ServiceAccountKey is the JSON representation of a service account key.
type ServiceAccountKey struct {
	// Name is the name (identity) of the service account key.
	Name string `json:"name"`

	// PGPKey is the armored PGP private key.
	PGPKey string `json:"pgp_key"`
}

// AuthInterceptorConfig defines Omni auth gRPC interceptors config.
type AuthInterceptorConfig struct {
	provider *client.KeyProvider

	// signer is the static signer to use if the config is sourced from the env variable ServiceAccountKeyEnvVar.
	signer message.Signer

	contextName string
	identity    string

	// envSource is set to true if the config is sourced from the environment variable ServiceAccountKeyEnvVar.
	envSource bool
}

// NewAuthInterceptorConfig creates new auth interceptor.
func NewAuthInterceptorConfig(contextName, identity, serviceAccountKey string) (*AuthInterceptorConfig, error) {
	if serviceAccountKey != "" {
		envIdentity, signer, err := parseServiceAccountKey(serviceAccountKey)
		if err != nil {
			return nil, err
		}

		return &AuthInterceptorConfig{
			envSource: true,
			identity:  envIdentity,
			signer:    signer,
		}, nil
	}

	return &AuthInterceptorConfig{
		envSource:   false,
		provider:    client.NewKeyProvider("omni/keys"),
		contextName: contextName,
		identity:    identity,
	}, nil
}

// Interceptor creates gRPC interceptor.
func (c *AuthInterceptorConfig) Interceptor() *interceptor.Signature {
	signerFunc := func(ctx context.Context, cc *grpc.ClientConn) (message.Signer, error) {
		// if the source is an environment variable, we can just return the static signer.
		if c.envSource {
			return c.signer, nil
		}

		return c.provider.ReadValidKey(c.contextName, c.identity)
	}

	// only attempt renews if config is not sourced from the env variables
	var renewSignerFunc interceptor.SignerFunc

	if !c.envSource {
		renewSignerFunc = func(ctx context.Context, cc *grpc.ClientConn) (message.Signer, error) {
			return c.authenticate(ctx, cc)
		}
	}

	authEnabledFunc := func(ctx context.Context, cc *grpc.ClientConn) (bool, error) {
		st := omni.NewClient(cc).State()
		confPtr := authres.NewAuthConfig().Metadata()

		// clear the outgoing metadata to prevent the request from being proxied to the Talos backend
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs())

		authConfig, err := safe.StateGet[*authres.Config](ctx, st, confPtr)
		if err != nil {
			return false, err
		}

		enabled := authConfig.TypedSpec().Value.GetAuth0().GetEnabled() || authConfig.TypedSpec().Value.GetWebauthn().GetEnabled()

		return enabled, nil
	}

	return interceptor.NewSignature(c.identity, signerFunc, renewSignerFunc, authEnabledFunc)
}

func (c *AuthInterceptorConfig) authenticate(ctx context.Context, cc *grpc.ClientConn) (*client.Key, error) {
	ctx = context.WithValue(ctx, interceptor.SkipInterceptorContextKey{}, struct{}{})

	authCli := authcli.NewClient(cc)

	err := c.provider.DeleteKey(c.contextName, c.identity)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	pgpKey, err := c.provider.GenerateKey(c.contextName, c.identity, version.Name+" "+version.Tag)
	if err != nil {
		return nil, err
	}

	publicKey, err := pgpKey.ArmorPublic()
	if err != nil {
		return nil, err
	}

	loginURL, err := authCli.RegisterPGPPublicKey(ctx, c.identity, []byte(publicKey))
	if err != nil {
		return nil, err
	}

	savePath, err := c.provider.WriteKey(pgpKey)
	if err != nil {
		return nil, err
	}

	printLoginDialog := func() {
		fmt.Fprintf(os.Stderr, "Please visit this page to authenticate with omni: %s\n", loginURL)
	}

	browserEnv := os.Getenv("BROWSER")
	if browserEnv == "echo" {
		printLoginDialog()
	} else {
		err = browser.OpenURL(loginURL)
		if err != nil {
			printLoginDialog()
		}
	}

	publicKeyID := pgpKey.Key.Fingerprint()

	err = authCli.AwaitPublicKeyConfirmation(ctx, publicKeyID)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Public key %s is now registered for user %s\n", publicKeyID, c.identity)

	fmt.Fprintf(os.Stderr, "PGP key saved to %s\n", savePath)

	return pgpKey, nil
}

func parseServiceAccountKey(value string) (string, message.Signer, error) {
	saKeyJSON, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", nil, err
	}

	var saKey ServiceAccountKey

	err = json.Unmarshal(saKeyJSON, &saKey)
	if err != nil {
		return "", nil, err
	}

	cryptoKey, err := pgpcrypto.NewKeyFromArmored(saKey.PGPKey)
	if err != nil {
		return "", nil, err
	}

	key, err := pgp.NewKey(cryptoKey)
	if err != nil {
		return "", nil, err
	}

	return saKey.Name, key, nil
}
