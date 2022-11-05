package connection

import (
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type Connection struct {
	grpcConn    *grpc.ClientConn
	tokenSource oauth2.TokenSource
	insecure    bool
}

func NewConnection(insecure bool, host string, audience string) (*Connection, error) {
	if insecure {
		return newConnectionInsecure(host)
	} else {
		return newConnectionSecure(host, audience)
	}
}

func newConnectionSecure(
	host string, audience string) (*Connection, error) {

	// build a secure Connection
	grpcConn, err := newGrpcConn(host)
	if err != nil {
		return nil, fmt.Errorf("newGrpcConn: %v", err)
	}

	tokenSource, err := newTokenSource(audience)
	if err != nil {
		return nil, fmt.Errorf("newTokenSource: %v", err)
	}

	return &Connection{
		grpcConn,
		tokenSource,
		false,
	}, nil
}

func newConnectionInsecure(
	host string) (*Connection, error) {

	// build an INSECURE connection
	grpcConn, err := newGrpcConnInsecure(host)
	if err != nil {
		return nil, fmt.Errorf("newGrpcConnInsecure: %v", err)
	}

	// does not bother creating a token source
	return &Connection{
		grpcConn,
		nil,
		true,
	}, nil
}

func (c *Connection) GrpcConn() *grpc.ClientConn {
	return c.grpcConn
}

func (c *Connection) Close() error {
	return c.grpcConn.Close()
}
