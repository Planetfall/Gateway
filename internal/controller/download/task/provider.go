package task

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
)

type Provider interface {
	NewClient() (Client, error)
}

type providerImpl struct {
}

func (p *providerImpl) NewClient() (Client, error) {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("provider.NewClient: %v", err)
	}

	return client, nil
}
