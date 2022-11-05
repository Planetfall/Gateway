package server

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	env "github.com/planetfall/gateway/pkg/environments"
)

var gcloudMapping = map[string](func(string, string) (*Gcloud, error)){
	env.Development: newGcloudDevelopment,
	env.Production:  newGcloudProduction,
}

// Google Cloud components wich behavior relies on the environment
type Gcloud struct {
	errorReporting *errorreporting.Client
	metadataClient *metadata.Client
}

func NewGcloud(
	env string, serviceName string, projectID string,
) (*Gcloud, error) {

	if handler, exists := gcloudMapping[env]; exists {
		return handler(serviceName, projectID)
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

func newGcloudDevelopment(_ string, _ string) (*Gcloud, error) {
	log.Printf("GCloud components disabled in %s\n", env.Development)
	return &Gcloud{
		errorReporting: nil,
		metadataClient: nil,
	}, nil
}

func newGcloudProduction(serviceName string, projectID string) (*Gcloud, error) {
	log.Printf("GCloud components enabled in %s, initializing...\n", env.Production)
	ctx := context.Background()

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
		metadataClient: nil,
		errorReporting: errorReporting,
	}, nil
}
