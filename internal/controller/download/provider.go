package download

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/planetfall/gateway/internal/controller/download/websocket"
)

// Provider which provides:
//   - the cloud task client
//   - the Pub/Sub subscriber
//   - the websocket upgrader
//   - the websocket store
type Provider interface {

	// Builds a new cloud task client.
	NewTaskClient(
		queuePath string, target string) (task.TaskClient, error)

	// Builds a new Pub/Sub client.
	NewSubscriber(
		onReceive func(ctx context.Context, message *pubsub.Message),
		projectID string, subscriptionID string,
		logger *log.Logger) (subscriber.Subscriber, error)

	// Builds a new websocket upgrader.
	// It also set the authorized origins for upgrades.
	NewWebsocket(
		origins []string) (websocket.Websocket, error)

	// Builds a new websocket store.
	NewWebsocketStore() websocket.Store
}

type providerImpl struct {
}

func (p *providerImpl) NewTaskClient(
	queuePath string, target string) (task.TaskClient, error) {

	// task client setup
	opt := task.TaskClientOptions{
		QueuePath: queuePath,
		Target:    target,
	}
	taskClient, err := task.NewTaskClient(opt)
	if err != nil {
		return nil, fmt.Errorf("task.NewTaskClient: %v", err)
	}

	return taskClient, nil
}

func (p *providerImpl) NewSubscriber(
	onReceive func(ctx context.Context, message *pubsub.Message),
	projectID string, subscriptionID string,
	logger *log.Logger) (subscriber.Subscriber, error) {

	opt := subscriber.SubscriberOptions{
		OnReceive:      onReceive,
		ProjectID:      projectID,
		SubscriptionID: subscriptionID,
		LoggerBase:     logger,
	}

	// pubsub client setup
	sub, err := subscriber.NewSubscriber(opt)
	if err != nil {
		return nil, fmt.Errorf("subscriber.NewSubscriber: %v", err)
	}

	return sub, nil
}

func (p *providerImpl) NewWebsocket(
	origins []string) (websocket.Websocket, error) {

	opt := websocket.WebsocketOptions{
		Origins: origins,
	}
	return websocket.NewWebsocket(opt)
}

func (p *providerImpl) NewWebsocketStore() websocket.Store {
	return websocket.NewStore()
}
