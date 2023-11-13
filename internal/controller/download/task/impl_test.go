package task_test

import (
	"fmt"
	"testing"

	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/planetfall/gateway/internal/controller/download/task/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask_withClientError(t *testing.T) {
	taskGiven := task.Task{
		JobKey: "key",
		Payload: task.Payload{
			Url:    "url",
			Artist: "artist",
			Album:  "album",
			Track:  "track",
		},
	}

	errorGiven := fmt.Errorf("failed to create task")
	clientGiven := mocks.NewClientMock().(*mocks.ClientMock)
	clientGiven.On("CreateTask").Return(&taskspb.Task{}, errorGiven)

	providerGiven := mocks.NewProviderMock(clientGiven).(*mocks.ProviderMock)
	queuePathGiven := "queue-path"
	targetGiven := "target"

	taskClient, err := task.NewTaskClient(task.TaskClientOptions{
		QueuePath: queuePathGiven,
		Target:    targetGiven,
		Provider:  providerGiven,
	})
	assert.Nil(t, err)

	_, err = taskClient.CreateTask(taskGiven)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cloudtasks.CreateTask")
}

func TestCreateTask(t *testing.T) {
	taskGiven := task.Task{
		JobKey: "key",
		Payload: task.Payload{
			Url:    "url",
			Artist: "artist",
			Album:  "album",
			Track:  "track",
		},
	}

	taskCreatedGiven := &taskspb.Task{
		Name: "task-name",
	}

	clientGiven := mocks.NewClientMock().(*mocks.ClientMock)
	clientGiven.On("CreateTask").Return(taskCreatedGiven, nil)

	providerGiven := mocks.NewProviderMock(clientGiven).(*mocks.ProviderMock)
	queuePathGiven := "queue-path"
	targetGiven := "target"

	taskClient, err := task.NewTaskClient(task.TaskClientOptions{
		QueuePath: queuePathGiven,
		Target:    targetGiven,
		Provider:  providerGiven,
	})
	assert.Nil(t, err)

	taskCreatedActual, err := taskClient.CreateTask(taskGiven)
	assert.Nil(t, err)

	assert.Equal(t, taskCreatedGiven.Name, taskCreatedActual.Name)
}
