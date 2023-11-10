// Package grpc provides a Connection implementation.
// It can be used to interact with GRPC services, hosted in secure or insecure
// environments.
//
// It aims to simplify the connection setup phase.
package grpc

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
)

// Connection is an helper, that holds the actual grpc.ClientConn
// connection.
// It encapsulate the authentication, and the setup of the transport
// credentials.
type Connection interface {
	Client() *grpc.ClientConn
	Close() error
	AuthenticateContext(context.Context) (context.Context, error)
}

type connectionImpl struct {
	client      *grpc.ClientConn
	tokenSource TokenSource
	insecure    bool
}

// Client allows read access on the grpc.ClientConn property.
func (c *connectionImpl) Client() *grpc.ClientConn {
	return c.client
}

// Close terminates the grpc.ClientConn connection
func (c *connectionImpl) Close() error {
	return c.client.Close()
}

// ConnectionOptions holds the parameters for the Connection builder
type ConnectionOptions struct {
	// Target build parameter
	Target string

	// Insecure builder parameter
	Insecure bool

	// Custom provider which provides a token source and a GRPC client
	// connection
	Provider Provider
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

	provider := getProvider(opt)

	client, err := provider.NewClient(opt.Target, opt.Insecure)
	if err != nil {
		return nil, fmt.Errorf("provider.NewClient: %v", err)
	}

	audience := buildAudienceFromTarget(opt.Target)
	tokenSource, err := provider.NewTokenSource(audience, opt.Insecure)
	if err != nil {
		return nil, fmt.Errorf("provider.NewTokenSource: %v", err)
	}

	return &connectionImpl{
		client:      client,
		tokenSource: tokenSource,
		insecure:    opt.Insecure,
	}, nil
}

// getProvider returns the provider given in options if not nil.
// Else, fallback to the default provider implementation.
func getProvider(opt ConnectionOptions) Provider {

	if opt.Provider != nil {
		return opt.Provider
	} else {
		return &providerImpl{}
	}
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
