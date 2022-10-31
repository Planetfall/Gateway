package connection

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// set up request context with authentication
// https://cloud.google.com/run/docs/triggering/grpc#request-auth
func newTokenSource(
	audience string) (oauth2.TokenSource, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create auth context with token
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		log.Printf("idtoken.NewTokenSource: %v", err)
		return nil, err
	}

	return tokenSource, nil
}

func (c *Connection) AuthenticateContext(ctx context.Context) (context.Context, error) {
	// if insecure set, does not authenticate
	if c.insecure == true {
		return ctx, nil
	}

	// else, add token to ctx
	token, err := c.tokenSource.Token()
	if err != nil {
		log.Printf("tokenSource.Token: %v", err)
		return nil, err
	}

	// set the token into a grpc context
	ctx = grpcMetadata.AppendToOutgoingContext(
		ctx, "authorization", "Bearer "+token.AccessToken)

	return ctx, nil
}
