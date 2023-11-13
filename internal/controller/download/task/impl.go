package task

import (
	"context"
	"encoding/json"
	"fmt"

	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

// taskClientImpl is the default implementation of the TaskClient
type taskClientImpl struct {
	client    Client
	queuePath string
	target    string
}

// TaskClientOptions are the options for the TaskClient builder.
type TaskClientOptions struct {
	// QueuePath locates where to push new tasks.
	QueuePath string

	// Target is the host that needs to be called by the task.
	Target string

	// Provider for the client
	Provider Provider
}

func (opt TaskClientOptions) getProvider() Provider {
	if opt.Provider == nil {
		return &providerImpl{}
	}

	return opt.Provider
}

// NewTaskClient is the builder for the TaskClient
func NewTaskClient(opt TaskClientOptions) (TaskClient, error) {

	provider := opt.getProvider()
	client, err := provider.NewClient()
	if err != nil {
		return nil, fmt.Errorf("task.NewTaskClient: %v", err)
	}
	return &taskClientImpl{
		client:    client,
		queuePath: opt.QueuePath,
		target:    opt.Target,
	}, nil
}

func (t *taskClientImpl) CreateTask(
	tPayload Task) (*taskspb.Task, error) {

	// json encode
	body, err := json.Marshal(&tPayload)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	req := t.newCreateTaskRequest(body)

	ctx := context.Background()
	createdTask, err := t.client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

func (t *taskClientImpl) newCreateTaskRequest(body []byte) *taskspb.CreateTaskRequest {

	return &taskspb.CreateTaskRequest{
		Parent: t.queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        t.target,
					Body:       body,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
		},
	}
}

func (t *taskClientImpl) Close() error {
	return t.client.Close()
}
