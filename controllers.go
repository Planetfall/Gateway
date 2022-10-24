package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/errorreporting"
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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

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

func formatError(err error, c *gin.Context, status int, message string) {
	log.Println(err)
	msg := errorMessage{
		Status:  status,
		Message: message,
	}
	c.JSON(status, msg)
	errorReporting.Report(errorreporting.Entry{
		Error: err,
	})
}

func badRequest(err error, c *gin.Context) {
	formatError(err, c,
		http.StatusBadRequest, "Wrong parameters supplied")
}

func internalError(err error, c *gin.Context) {
	formatError(err, c,
		http.StatusInternalServerError, "Something went wrong on my side")
}
