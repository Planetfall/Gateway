package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	env "github.com/planetfall/gateway/pkg/environments"
)

var gcloudMapping = map[string](func(string) (*Gcloud, error)){
	env.Development: newGcloudDevelopment,
	env.Production:  newGcloudProduction,
}

// Google Cloud components wich behavior relies on the environment
type Gcloud struct {
	errorReporting *errorreporting.Client
	metadataClient *metadata.Client
}

func NewGcloud(env string, serviceName string) (*Gcloud, error) {
	if handler, exists := gcloudMapping[env]; exists {
		return handler(serviceName)
	}

	return nil, fmt.Errorf("could not initialize gcloud components with env %v", env)
}

func (g *Gcloud) ErrorReport(err error) {
	if g.errorReporting != nil {
		g.errorReporting.Report(errorreporting.Entry{
			Error: err,
		})
	}
}

func newGcloudDevelopment(_ string) (*Gcloud, error) {
	log.Printf("GCloud components disabled in %s\n", env.Development)
	return &Gcloud{
		errorReporting: nil,
		metadataClient: nil,
	}, nil
}

func newGcloudProduction(serviceName string) (*Gcloud, error) {
	log.Printf("GCloud components enabled in %s\n", env.Production)
	ctx := context.Background()

	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		return nil, fmt.Errorf("metadataClient.ProjectID: %v", err)
	}

	errorReporting, err := errorreporting.NewClient(
		ctx, projectID,
		errorreporting.Config{
			ServiceName: serviceName,
			OnError: func(err error) {
				log.Printf("Could not log error: %v\n", err)
			}},
	)
	if err != nil {
		return nil, fmt.Errorf("errorreporting.NewClient: %v", err)
	}

	return &Gcloud{
		metadataClient: metadataClient,
		errorReporting: errorReporting,
	}, nil
}
