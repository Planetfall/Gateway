package mocks

import (
	"context"

	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/googleapis/gax-go"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() task.Client {
	return &ClientMock{}
}

func (m *ClientMock) CreateTask(ctx context.Context,
	req *taskspb.CreateTaskRequest, opts ...gax.CallOption) (*taskspb.Task, error) {

	args := m.Called()
	return args.Get(0).(*taskspb.Task), args.Error(1)
}

func (m *ClientMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
