package task

import (
	"context"
	"encoding/json"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type taskClientImpl struct {
	client    *cloudtasks.Client
	queuePath string
	target    string
}

func (t *taskClientImpl) CreateTask(
	tPayload TaskPayload) (*taskspb.Task, error) {

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
