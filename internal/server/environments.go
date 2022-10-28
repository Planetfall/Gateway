package server

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
)

const (
	Development string = "development"
	Production  string = "production"
)

func newServerDevelopment(
	serviceName string, port string, conns connections, cls *clients,
) (*Server, error) {
	return &Server{
		port:           port,
		router:         nil,
		errorReporting: nil,
		metadataClient: nil,
		serviceName:    serviceName,
		conns:          conns,
		cls:            cls,
	}, nil
}

func newServerProduction(
	serviceName string, port string, conns connections, cls *clients,
) (*Server, error) {

	ctx := context.Background()

	log.Println("initializing metadata client...")
	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		log.Printf("metadata.ProjectID: %v\n", err)
		return nil, err
	}

	log.Println("initializing error reporting...")
	errorReporting, err := errorreporting.NewClient(ctx, projectID,
		errorreporting.Config{
			ServiceName: serviceName,
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			}},
	)
	if err != nil {
		log.Printf("errorreporting.NewClient: %v\n", err)
		return nil, err
	}

	return &Server{
		port:           port,
		router:         nil,
		errorReporting: errorReporting,
		metadataClient: metadataClient,
		serviceName:    serviceName,
		conns:          conns,
		cls:            cls,
	}, nil
}
