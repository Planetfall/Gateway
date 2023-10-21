// Package download aims to interact with the download job using Cloud and
// websocket features.
package download

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gorilla/websocket"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"

	"github.com/planetfall/gateway/internal/controller"
)

// DownloadController is used to interact with the youtube download job
type DownloadController struct {
	// Reference to the base controller type
	controller.Controller

	// The path used when created new tasks in Cloud Tasks
	queuePath string

	// The client to push new tasks in Cloud Tasks
	tasks *cloudtasks.Client

	// The subscriber helper to pull Pub/Sub messages
	sub *subscriber.Subscriber

	// The upgrader to upgrade HTTP request to websocket
	upgrader websocket.Upgrader

	// The list of allowed origins that can interact with this controller using
	// websockets
	origins []string
}

// DownloadControllerOptions holds the parameters for the DownloadController
// builder
type DownloadControllerOptions struct {
	// The [controller] builder parameters
	ControllerOptions controller.ControllerOptions

	// The project ID in Google Cloud
	ProjectID string

	// The location ID configured in Cloud Tasks
	LocationID string

	// The queue ID where to push new tasks in Cloud Tasks
	QueueID string

	// The subscribe ID to use when pulling messages from Pub/Sub
	SubscriptionID string

	// The list of allowed origins that can interact with this controller using
	// websockets
	Origins []string
}

// NewDownloadController builds a new DownloadController.
// It setup the task client and the Pub/Sub helper.
func NewDownloadController(opt DownloadControllerOptions) (*DownloadController, error) {

	// initialize the base type
	ctrl := controller.NewController(opt.ControllerOptions)

	// builds the task queue path
	queuePath := fmt.Sprintf(
		"projects/%s/locations/%s/queues/%s",
		opt.ProjectID, opt.LocationID, opt.QueueID,
	)

	// task client setup
	ctx := context.Background()
	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.NewClient: %v", err)
	}

	// pubsub client setup
	sub, err := subscriber.NewSubscriber(
		opt.ProjectID, opt.SubscriptionID, ctrl.Logger)
	if err != nil {
		return nil, fmt.Errorf("subscriber.NewSubscriber: %v", err)
	}
	// starts listening for pubsub messages
	go func() {
		if err := sub.Listen(); err != nil {
			ctrl.Logger.Println(fmt.Errorf("subscriber.Listen: %v", err))
		}
	}()

	// init the controller
	downloadCtrl := &DownloadController{
		Controller: ctrl,

		queuePath: queuePath,
		tasks:     taskClient,
		sub:       sub,
		origins:   opt.Origins,
	}

	// setup the websocket upgrader
	downloadCtrl.upgrader = websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin:     downloadCtrl.checkOrigins,
	}

	return downloadCtrl, nil
}

// Close closes the task client and the Pub/Sub helper
func (c *DownloadController) Close() error {
	if err := c.tasks.Close(); err != nil {
		return fmt.Errorf("cloudtasks.Close: %v", err)
	}

	c.sub.Close()

	return nil
}
