package mocks

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/planetfall/gateway/internal/controller/download"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/planetfall/gateway/internal/controller/download/websocket"
	"github.com/stretchr/testify/mock"
)

func NewProviderMock(
	taskClient task.TaskClient,
	subscriber subscriber.Subscriber,
	websocket websocket.Websocket,
	store websocket.Store) download.Provider {

	p := &ProviderMock{}
	p.On("NewTaskClient").Return(taskClient, nil)
	p.On("NewSubscriber").Return(subscriber, nil)
	p.On("NewWebsocket").Return(websocket, nil)
	p.On("NewWebsocketStore").Return(store)
	return p
}

type ProviderMock struct {
	mock.Mock
}

func (m *ProviderMock) NewTaskClient(
	queuePath string, target string) (task.TaskClient, error) {

	args := m.Called()
	return args.Get(0).(task.TaskClient), args.Error(1)
}

func (m *ProviderMock) NewSubscriber(
	onReceive func(ctx context.Context, message *pubsub.Message),
	projectID string, subscriptionID string,
	logger *log.Logger) (subscriber.Subscriber, error) {

	args := m.Called()
	return args.Get(0).(subscriber.Subscriber), args.Error(1)
}

func (m *ProviderMock) NewWebsocket(
	origins []string) (websocket.Websocket, error) {

	args := m.Called()
	return args.Get(0).(websocket.Websocket), args.Error(1)
}

func (m *ProviderMock) NewWebsocketStore() websocket.Store {
	args := m.Called()
	return args.Get(0).(websocket.Store)
}
