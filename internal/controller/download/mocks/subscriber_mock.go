package mocks

import (
	"cloud.google.com/go/pubsub"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/stretchr/testify/mock"
)

func NewSubscriberMock() subscriber.Subscriber {
	return &SubscriberMock{}
}

type SubscriberMock struct {
	mock.Mock
}

func (m *SubscriberMock) Listen() error {
	args := m.Called()
	return args.Error(0)
}

func (m *SubscriberMock) Close() {
	m.Called()
}

func (m *SubscriberMock) NewJobStatus(
	message *pubsub.Message) (*subscriber.JobStatus, error) {

	args := m.Called(message)
	return args.Get(0).(*subscriber.JobStatus), args.Error(1)
}
