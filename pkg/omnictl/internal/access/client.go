// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package access

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/cosi-project/runtime/pkg/state"

	"github.com/siderolabs/omni-client/pkg/client"
	"github.com/siderolabs/omni-client/pkg/omni/resources"
	"github.com/siderolabs/omni-client/pkg/omni/resources/system"
	"github.com/siderolabs/omni-client/pkg/omnictl/config"
	"github.com/siderolabs/omni-client/pkg/version"
)

const (
	// EndpointEnvVar is the name of the environment variable that contains the Omni endpoint.
	EndpointEnvVar = "OMNI_ENDPOINT"
	// ServiceAccountKeyEnvVar is the name of the environment variable that contains the base64-encoded service account key JSON.
	ServiceAccountKeyEnvVar = "OMNI_SERVICE_ACCOUNT_KEY"
)

type clientOptions struct {
	skipAuth bool
}

// ClientOption is a functional option for the client.
type ClientOption func(*clientOptions)

// WithSkipAuth configures the client to skip the authentication interception.
func WithSkipAuth(skipAuth bool) ClientOption {
	return func(o *clientOptions) {
		o.skipAuth = skipAuth
	}
}

// WithClient initializes the Omni API client.
func WithClient(f func(ctx context.Context, client *client.Client) error, clientOpts ...ClientOption) error {
	_, err := config.Init(CmdFlags.Omniconfig, true)
	if err != nil {
		return err
	}

	cliOpts := clientOptions{}

	for _, opt := range clientOpts {
		opt(&cliOpts)
	}

	return WithContext(func(ctx context.Context) error {
		var opts []client.Option

		conf, err := config.Current()
		if err != nil {
			return err
		}

		contextName := conf.Context
		if CmdFlags.Context != "" {
			contextName = CmdFlags.Context
		}

		configCtx, err := conf.GetContext(CmdFlags.Context)
		if err != nil {
			return err
		}

		if configCtx.Auth.Basic != "" {
			_, _, found := strings.Cut(configCtx.Auth.Basic, ":")
			if !found {
				return fmt.Errorf("auth-basic should be in the format of <username>:<password>")
			}

			opts = append(opts, client.WithBasicAuth(configCtx.Auth.Basic))
		}

		serviceAccountKey := os.Getenv(ServiceAccountKeyEnvVar)
		if serviceAccountKey != "" {
			opts = append(opts, client.WithServiceAccount(contextName, serviceAccountKey))
		} else {
			opts = append(opts, client.WithUserAccount(contextName, configCtx.Auth.SideroV1.Identity))
		}

		if configCtx.URL == config.PlaceholderURL {
			return fmt.Errorf("context %q has not been configured, you will need to set it manually", contextName)
		}

		url := configCtx.URL

		endpointEnv := os.Getenv(EndpointEnvVar)
		if endpointEnv != "" {
			url = endpointEnv
		}

		client, err := client.New(ctx, url, opts...)
		if err != nil {
			return err
		}

		if !cliOpts.skipAuth {
			// bootstrap the client, and perform auth/re-auth if needed via the unary call
			// stream interceptor can't catch the auth error, as it comes async
			_, err = client.Omni().State().Get(ctx, resource.NewMetadata(resources.EphemeralNamespace, system.SysVersionType, system.SysVersionID, resource.VersionUndefined))
			if err != nil {
				return err
			}
		}

		if err = checkVersion(ctx, client.Omni().State()); err != nil {
			return err
		}

		return f(ctx, client)
	})
}

func checkVersion(ctx context.Context, state state.State) error {
	if version.Tag == "" && !version.SuppressVersionWarning {
		fmt.Println(`[WARN] github.com/siderolabs/omni-client/pkg/version.Tag is not set, client-server version validation is disabled.
If you want to disable this warning set github.com/siderolabs/omni-client/pkg/version.SuppressVersionWarning to true.`)

		return nil
	}

	sysversion, err := safe.StateGet[*system.SysVersion](ctx, state, system.NewSysVersion(resources.EphemeralNamespace, system.SysVersionID).Metadata())
	if err != nil {
		return err
	}

	backendVersion, err := semver.Parse(strings.TrimLeft(sysversion.TypedSpec().Value.BackendVersion, "v"))
	if err != nil {
		return err
	}

	clientVersion, err := semver.Parse(strings.TrimLeft(version.Tag, "v"))
	if err != nil {
		return err
	}

	if backendVersion.Major != clientVersion.Major || backendVersion.Minor != clientVersion.Minor {
		return fmt.Errorf("client version mismatch: backend version %s, client version %s", sysversion.TypedSpec().Value.BackendVersion, version.Tag)
	}

	return nil
}
