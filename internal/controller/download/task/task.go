package task

import (
	"context"

	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/googleapis/gax-go"
)

// Task is the payload sent to the download job.
type Task struct {
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

type Client interface {
	CreateTask(ctx context.Context, req *taskspb.CreateTaskRequest, opts ...gax.CallOption) (*taskspb.Task, error)
	Close() error
}

// TaskClient is responsible for creating tasks in Cloud Task.
type TaskClient interface {

	// CreateTask creates a new task from the given payload.
	// It returns the created task.
	CreateTask(tPayload Task) (*taskspb.Task, error)

	// Close closes the client
	Close() error
}
