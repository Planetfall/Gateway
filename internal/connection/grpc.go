package connection

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
