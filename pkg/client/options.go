// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package client provides Omni API client.
package client

import (
	"google.golang.org/grpc"

	"github.com/siderolabs/omni-client/pkg/access"
)

// Option is the function that generates gRPC dial options.
type Option func() ([]grpc.DialOption, error)

// WithBasicAuth creates the client with basic auth.
func WithBasicAuth(auth string) Option {
	return func() ([]grpc.DialOption, error) {
		return []grpc.DialOption{grpc.WithPerRPCCredentials(BasicAuth{
			auth: auth,
		})}, nil
	}
}

// WithServiceAccount creates the client for a context with the given service account key.
func WithServiceAccount(contextName, key string) Option {
	return func() ([]grpc.DialOption, error) {
		interceptorConfig, err := access.NewAuthInterceptorConfig(contextName, "", key)
		if err != nil {
			return nil, err
		}

		authInterceptor := interceptorConfig.Interceptor()

		return []grpc.DialOption{
			grpc.WithUnaryInterceptor(authInterceptor.Unary()),
			grpc.WithStreamInterceptor(authInterceptor.Stream()),
		}, nil
	}
}

// WithUserAccount used for accessing Omni by a human.
func WithUserAccount(contextName, identity string) Option {
	return func() ([]grpc.DialOption, error) {
		interceptorConfig, err := access.NewAuthInterceptorConfig(contextName, identity, "")
		if err != nil {
			return nil, err
		}

		authInterceptor := interceptorConfig.Interceptor()

		return []grpc.DialOption{
			grpc.WithUnaryInterceptor(authInterceptor.Unary()),
			grpc.WithStreamInterceptor(authInterceptor.Stream()),
		}, nil
	}
}

// WithGrpcOpts creates the client with basic auth.
func WithGrpcOpts(opts ...grpc.DialOption) Option {
	return func() ([]grpc.DialOption, error) {
		return opts, nil
	}
}
