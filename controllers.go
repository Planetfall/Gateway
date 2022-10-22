package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// init grpc conn for a GCloud microservice
// https://cloud.google.com/run/docs/triggering/grpc#connect
func getGrpcConn(host string, audience string) (*grpc.ClientConn, error) {
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
func getAuthenticatedCtx(conn *grpc.ClientConn) (context.Context, error) {

	// temp ctx to retrieve an auth token
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create auth context with token
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		log.Printf("idtoken.NewTokenSource: %v", err)
		return nil, err
	}

	token, err := tokenSource.Token()
	if err != nil {
		log.Printf("tokenSource.Token: %v", err)
		return nil, err
	}

	// set the token into a grpc context
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	return ctx, nil
}

// errors
type errorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func badRequest(c *gin.Context) {
	status := http.StatusBadRequest
	msg := errorMessage{
		Status:  status,
		Message: "Wrong parameters supplied",
	}
	c.JSON(status, &msg)
}

func internalError(c *gin.Context) {
	status := http.StatusInternalServerError
	msg := errorMessage{
		Status:  status,
		Message: "Something went wrong on my side",
	}
	c.JSON(status, &msg)
}
