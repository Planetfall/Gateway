package server

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/errorreporting"
)

// Google Cloud components
type Gcloud struct {
	errorReporting *errorreporting.Client
}

func NewGcloud(serviceName string, projectID string) (*Gcloud, error) {
	log.Printf("initializing GCloud components")
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
		errorReporting: errorReporting,
	}, nil
}

func (g *Gcloud) Close() error {
	if err := g.errorReporting.Close(); err != nil {
		return fmt.Errorf("errorReporting.Close: %v")
	}

	return nil
}

func (g *Gcloud) ErrorReport(err error) {
	if g.errorReporting != nil {
		g.errorReporting.Report(errorreporting.Entry{
			Error: err,
		})
	}
}
