package task

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

// TaskPayload is the payload sent to the download job.
type TaskPayload struct {
	JobKey string `json:"job_key"`
	Payload
}

// Payload contains the needed fields to perform the download.
type Payload struct {
	Url    string `json:"url"`    // the youtube url to use for download
	Artist string `json:"artist"` // the artist
	Album  string `json:"album"`  // the album
	Track  string `json:"track"`  // the music track
}

// TaskClient is responsible for creating tasks in Cloud Task.
type TaskClient interface {
	// CreateTask creates a new task from the given payload.
	// It returns the created task.
	CreateTask(tPayload TaskPayload) (*taskspb.Task, error)

	// Close closes the client
	Close() error
}

// TaskClientOptions are the options for the TaskClient builder.
type TaskClientOptions struct {
	// QueuePath locates where to push new tasks.
	QueuePath string

	// Target is the host that needs to be called by the task.
	Target string
}

// NewTaskClient is the builder for the TaskClient
func NewTaskClient(opt TaskClientOptions) (TaskClient, error) {

	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &taskClientImpl{
		client:    client,
		queuePath: opt.QueuePath,
		target:    opt.Target,
	}, nil
}
