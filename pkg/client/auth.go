// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package client

import (
	"context"
	"encoding/base64"
)

// BasicAuth adds basic auth for each gRPC request.
type BasicAuth struct {
	auth string
}

// GetRequestMetadata implements credentials.PerGRPCCredentials.
func (c BasicAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	enc := base64.StdEncoding.EncodeToString([]byte(c.auth))

	return map[string]string{
		"Authorization": "Basic " + enc,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials.
func (c BasicAuth) RequireTransportSecurity() bool {
	return true
}
