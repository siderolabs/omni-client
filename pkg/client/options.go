// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package client provides Omni API client.
package client

import (
	"context"

	"github.com/cosi-project/runtime/pkg/safe"
	"github.com/siderolabs/go-api-signature/pkg/client/interceptor"
	"github.com/siderolabs/go-api-signature/pkg/pgp/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/siderolabs/omni-client/pkg/client/omni"
	authres "github.com/siderolabs/omni-client/pkg/omni/resources/auth"
	"github.com/siderolabs/omni-client/pkg/version"
)

// Options is the options for the client.
type Options struct {
	BasicAuth string

	AuthInterceptor *interceptor.Interceptor

	AdditionalGRPCDialOptions []grpc.DialOption
}

// Option is a functional option for the client.
type Option func(*Options)

// WithBasicAuth creates the client with basic auth.
func WithBasicAuth(auth string) Option {
	return func(options *Options) {
		options.BasicAuth = auth
	}
}

// WithServiceAccount creates the client authenticating with the given service account.
func WithServiceAccount(serviceAccountBase64 string) Option {
	return func(options *Options) {
		options.AuthInterceptor = signatureAuthInterceptor("", "", serviceAccountBase64)
	}
}

// WithUserAccount is used for accessing Omni by a human.
func WithUserAccount(contextName, identity string) Option {
	return func(options *Options) {
		options.AuthInterceptor = signatureAuthInterceptor(contextName, identity, "")
	}
}

func signatureAuthInterceptor(contextName, identity, serviceAccountBase64 string) *interceptor.Interceptor {
	return interceptor.New(interceptor.Options{
		AuthEnabledFunc: func(ctx context.Context, cc *grpc.ClientConn) (bool, error) {
			st := omni.NewClient(cc).State()
			confPtr := authres.NewAuthConfig().Metadata()

			// clear the outgoing metadata to prevent the request from being proxied to the Talos backend
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs())

			authConfig, err := safe.StateGet[*authres.Config](ctx, st, confPtr)
			if err != nil {
				return false, err
			}

			enabled := authres.Enabled(authConfig)

			return enabled, nil
		},
		UserKeyProvider:      client.NewKeyProvider("omni/keys"),
		ContextName:          contextName,
		Identity:             identity,
		ClientName:           version.Name + " " + version.Tag,
		ServiceAccountBase64: serviceAccountBase64,
	})
}

// WithGrpcOpts adds additional gRPC dial options to the client.
func WithGrpcOpts(opts ...grpc.DialOption) Option {
	return func(options *Options) {
		options.AdditionalGRPCDialOptions = append(options.AdditionalGRPCDialOptions, opts...)
	}
}
