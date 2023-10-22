package grpc

import (
	"context"
	"fmt"

	grpcMetadata "google.golang.org/grpc/metadata"
)

// AuthenticateContext enrich an input context with an authentication token.
// It retrieve this token using the connection configured token source.
// If this source is not set, it returns the provided context unchanged.
// This is reused from the [Cloud Run] documentation
//
// [Cloud Run]: https://cloud.google.com/run/docs/triggering/grpc#request-auth
func (c *connectionImpl) AuthenticateContext(
	ctx context.Context) (context.Context, error) {

	// if tokenSource unset, not able to provide a token
	if c.tokenSource == nil {
		return ctx, nil
	}

	// else, add token to ctx
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("tokenSource.Token: %v", err)
	}

	// set the token into a grpc context
	ctx = grpcMetadata.AppendToOutgoingContext(
		ctx, "authorization", "Bearer "+token.AccessToken)

	return ctx, nil
}
