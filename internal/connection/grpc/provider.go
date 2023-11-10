package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Provider is responsible to provide new instances of:
//   - a new GRPC client connection
//   - a new token source
type Provider interface {

	// NewTokenSource builds a new token source instance.
	NewTokenSource(audience string, insecure bool) (TokenSource, error)

	// NewClient builds a new GRPC client connection
	NewClient(target string, insecure bool) (*grpc.ClientConn, error)
}

// The default provider implementation
type providerImpl struct {
}

// NewClient creates a new GRPC connection.
// It retrieves transport credentials if insecure is false.
func (p *providerImpl) NewClient(
	target string, insecure bool) (*grpc.ClientConn, error) {

	creds, err := p.newCredentials(insecure)
	if err != nil {
		return nil, fmt.Errorf("getCredentials: %v", err)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}
	client, err := grpc.Dial(target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %v", err)
	}

	return client, nil
}

// newCredentials provides transport credentials according to the isInsecure
// argument:
//   - if true, it use the [insecure] library to provide empty credentials
//   - else, it will provide valid TLS transport credentials
func (p *providerImpl) newCredentials(isInsecure bool) (credentials.TransportCredentials, error) {
	if isInsecure {
		return insecure.NewCredentials(), nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("x509.SystemCertPool: %v", err)
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})

	return creds, nil
}

// NewTokenSource builds a new source which provides authentication tokens.
// The audience is the target service host that we need authentication for.
// This is reused from the [Cloud Run] documentation
//
// [Cloud Run]: https://cloud.google.com/run/docs/triggering/grpc#request-auth
func (p *providerImpl) NewTokenSource(
	audience string, insecure bool) (TokenSource, error) {

	if insecure {
		return nil, nil
	}

	ctx := context.Background()

	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewTokenSource: %v", err)
	}

	return tokenSource, nil
}
