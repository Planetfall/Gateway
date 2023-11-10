package mocks

import (
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/stretchr/testify/mock"
)

func NewTaskClientMock() task.TaskClient {
	return &TaskClientMock{}
}

type TaskClientMock struct {
	mock.Mock
}

func (m *TaskClientMock) CreateTask(
	tPayload task.TaskPayload) (*taskspb.Task, error) {

	args := m.Called()
	return args.Get(0).(*taskspb.Task), args.Error(1)
}

func (m *TaskClientMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
