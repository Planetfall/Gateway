package download

import (
	"context"
	"encoding/json"
	"fmt"

	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

func (c *DownloadController) createTask(
	dPayload taskDownloadPayload) (*taskspb.Task, error) {

	body, err := json.Marshal(&dPayload)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	req := c.newCreateTaskRequest(body)

	ctx := context.Background()
	createdTask, err := c.tasks.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

func (c *DownloadController) newCreateTaskRequest(body []byte) *taskspb.CreateTaskRequest {

	return &taskspb.CreateTaskRequest{
		Parent: c.queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        c.Target,
					Body:       body,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
				},
			},
		},
	}
}
