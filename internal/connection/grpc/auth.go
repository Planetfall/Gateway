package grpc

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// TokenSource is responsible for providing OAuth tokens
type TokenSource interface {

	// Token provides a valid OAuth token
	Token() (*oauth2.Token, error)
}

// AuthenticateContext enrich an input context with an authentication token.
// It retrieve this token using the connection configured token source.
// If insecure is explicited provided, the given context is returned unchanged.
// This is reused from the [Cloud Run] documentation
//
// [Cloud Run]: https://cloud.google.com/run/docs/triggering/grpc#request-auth
func (c *connectionImpl) AuthenticateContext(
	ctx context.Context) (context.Context, error) {

	// if tokenSource unset, not able to provide a token
	if c.insecure {
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
