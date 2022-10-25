package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
)

type connection struct {
	grpcConn    *grpc.ClientConn
	tokenSource oauth2.TokenSource
}

func newConnection(host string, audience string) (*connection, error) {
	c := &connection{}
	grpcConn, err := c.newGrpcConn(host)
	if err != nil {
		log.Printf("failed to create connection for %s\n", host)
		return nil, err
	}

	tokenSource, err := c.newTokenSource(audience)
	if err != nil {
		log.Printf("failed to get token token source for %s\n", audience)
		return nil, err
	}
	c.grpcConn = grpcConn
	c.tokenSource = tokenSource

	return c, nil
}

// init grpc conn for a GCloud microservice
// https://cloud.google.com/run/docs/triggering/grpc#connect
func (c *connection) newGrpcConn(host string) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Println("failed getting certs")
		return nil, err
	}

	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("failed creating the GRPC connection with %s\n", host)
		return nil, err
	}

	return conn, nil
}

// set up request context with authentication
// https://cloud.google.com/run/docs/triggering/grpc#request-auth
func (c *connection) newTokenSource(
	audience string) (oauth2.TokenSource, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// create auth context with token
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		log.Printf("idtoken.NewTokenSource: %v", err)
		return nil, err
	}

	return tokenSource, nil
}

func (c *connection) getAuthenticatedCtx(ctx context.Context) (context.Context, error) {

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
