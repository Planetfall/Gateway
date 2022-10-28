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
	insecure    bool
}

func newConnection(
	host string, audience string) (*connection, error) {

	// build a secure connection
	grpcConn, err := newGrpcConn(host)
	if err != nil {
		log.Printf("failed to create connection to %s\n", host)
		return nil, err
	}

	tokenSource, err := newTokenSource(audience)
	if err != nil {
		log.Printf("failed to get token token source for %s\n", audience)
		return nil, err
	}

	return &connection{
		grpcConn,
		tokenSource,
		false,
	}, nil
}

func newConnectionInsecure(host string) (*connection, error) {

	// build an INSECURE connection
	grpcConn, err := newGrpcConnInsecure(host)
	if err != nil {
		log.Printf("failed to create INSECURE connection to %s\n", host)
		return nil, err
	}

	return &connection{
		grpcConn,
		nil,
		true,
	}, nil
}

// init grpc conn for a GCloud microservice
// https://cloud.google.com/run/docs/triggering/grpc#connect
func newGrpcConnInsecure(
	host string) (*grpc.ClientConn, error) {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("grpc.Dial: %v\n", err)
		return nil, err
	}

	return conn, nil
}

// init grpc conn for a GCloud microservice using TLS
// https://cloud.google.com/run/docs/triggering/grpc#connect
func newGrpcConn(host string) (*grpc.ClientConn, error) {

	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("x509.SystemCertPool: %v\n", err)
		return nil, err
	}

	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Printf("grpc.Dial: %v\n", err)
		return nil, err
	}

	return conn, nil
}

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

func (c *connection) getAuthenticatedCtx(ctx context.Context) (context.Context, error) {

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

func (c *connection) Close() error {
	return c.grpcConn.Close()
}
