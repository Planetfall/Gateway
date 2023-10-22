// Package grpc provides a Connection implementation.
// It can be used to interact with GRPC services, hosted in secure or insecure
// environments.
//
// It aims to simplify the connection setup phase.
package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Connection is an helper, that holds the actual grpc.ClientConn
// connection.
// It encapsulate the authentication, and the setup of the transport
// credentials.
type Connection interface {
	GrpcConn() *grpc.ClientConn
	Close() error
	AuthenticateContext(context.Context) (context.Context, error)
}

type connectionImpl struct {
	grpcConn    *grpc.ClientConn
	tokenSource oauth2.TokenSource
}

// GrpcConn allows read access on the grpc.ClientConn property.
func (c *connectionImpl) GrpcConn() *grpc.ClientConn {
	return c.grpcConn
}

// Close terminates the grpc.ClientConn connection
func (c *connectionImpl) Close() error {
	return c.grpcConn.Close()
}

// ConnectionOptions holds the parameters for the Connection builder
type ConnectionOptions struct {
	// Target build parameter
	Target string

	// Insecure builder parameter
	Insecure bool
}

// NewConnection builds a new GRPC connection object.
// It holds the actual grpc.ClientConn used to interact with a GRPC service.
// It also provides an helper method to authenticate a context.
//
// The isInsecure parameter is used to set on/off the security context:
//   - if true, the token source is unset, and the transport credentials are
//     empty.
//   - else, the token source is configured, and and the transport credentials
//     use TLS.
//
// The host is used to setup the grpc.ClientConn.
func NewConnection(opt ConnectionOptions) (Connection, error) {
	creds, err := newCredentials(opt.Insecure)
	if err != nil {
		return nil, fmt.Errorf("getCredentials: %v", err)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}
	grpcConn, err := grpc.Dial(opt.Target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial: %v", err)
	}

	audience := buildAudienceFromTarget(opt.Target)
	var tokenSource oauth2.TokenSource = nil
	if !opt.Insecure {
		tokenSource, err = newTokenSource(audience)
		if err != nil {
			return nil, fmt.Errorf("newTokenSource: %v", err)
		}
	}

	return &connectionImpl{
		grpcConn,
		tokenSource,
	}, nil
}

// newCredentials provides transport credentials according to the isInsecure
// argument:
//   - if true, it use the [insecure] library to provide empty credentials
//   - else, it will provide valid TLS transport credentials
func newCredentials(isInsecure bool) (credentials.TransportCredentials, error) {
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

// buildAudienceFromTarget builds the audience using the configured target.
// As audience is supposed to be used alonside a token source, it assumes
// that the target is running with TLS. It prepend the output audience with the
// HTTPS scheme.
//
// For example, the following targets:
//   - music-researcher.run.app:443
//   - locahost:8080
//
// Should have the following audiences:
//   - https://music-researcher.run.app
//   - https://localhost:8080
//
// It is generated through this method and not configured to avoid duplication
// of the host.
func buildAudienceFromTarget(target string) string {
	// removes the port part if any
	host := strings.Split(target, ":")[0]

	// audience is only used in secured/TLS environment
	// we can safely assume that only HTTPS scheme is needed
	audience := fmt.Sprintf("https://%s", host)

	return audience
}

// newTokenSource builds a new source which provides authentication tokens.
// The audience is the target service host that we need authentication for.
// This is reused from the [Cloud Run] documentation
//
// [Cloud Run]: https://cloud.google.com/run/docs/triggering/grpc#request-auth
func newTokenSource(audience string) (oauth2.TokenSource, error) {

	ctx := context.Background()

	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("idtoken.NewTokenSource: %v", err)
	}

	return tokenSource, nil
}
