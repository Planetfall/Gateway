// Package download aims to interact with the download job using Cloud and
// websocket features.
package download

import (
	"fmt"

	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/planetfall/gateway/internal/controller/download/websocket"
)

// DownloadController is used to interact with the youtube download job
type DownloadController struct {
	// Reference to the base controller type
	controller.Controller

	// The path used when created new tasks in Cloud Tasks
	queuePath string

	// The client to push new taskClient in Cloud Tasks
	taskClient task.TaskClient

	// The subscriber helper to pull Pub/Sub messages
	sub subscriber.Subscriber

	// The upgrader to upgrade HTTP request to websocket
	websocket websocket.Websocket

	// The store which holds the active websockets and their job keys.
	// It is used to retrieve the concerned websocket when a message is
	// received.
	websocketStore websocket.Store
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

	// Custom provider
	Provider Provider
}

func (opt DownloadControllerOptions) getProvider() Provider {
	if opt.Provider == nil {
		return &providerImpl{}
	}

	return opt.Provider
}

// NewDownloadController builds a new DownloadController.
// It setup the task client and the Pub/Sub helper.
func NewDownloadController(
	opt DownloadControllerOptions) (*DownloadController, error) {

	// initialize the base type
	ctrl := controller.NewController(opt.ControllerOptions)

	// setup the download controller
	downloadCtrl := &DownloadController{
		Controller: ctrl,
	}

	// retrieve the provider
	provider := opt.getProvider()

	// setup the task client

	// builds the task queue path
	queuePath := fmt.Sprintf(
		"projects/%s/locations/%s/queues/%s",
		opt.ProjectID, opt.LocationID, opt.QueueID,
	)
	taskClient, err := provider.NewTaskClient(queuePath, ctrl.Target)
	if err != nil {
		return nil, fmt.Errorf("provider.NewTaskClient: %v", err)
	}
	downloadCtrl.queuePath = queuePath
	downloadCtrl.taskClient = taskClient

	// setup the subscriber
	sub, err := provider.NewSubscriber(
		downloadCtrl.OnReceive,
		opt.ProjectID,
		opt.SubscriptionID,
		opt.ControllerOptions.Logger,
	)
	if err != nil {
		return nil, fmt.Errorf("provider.NewSubscriber: %v", err)
	}

	// starts listening for pubsub messages
	go func() {
		if err := sub.Listen(); err != nil {
			ctrl.Logger.Println(fmt.Errorf("subscriber.Listen: %v", err))
		}
	}()
	downloadCtrl.sub = sub

	// setup the websocket store
	store := provider.NewWebsocketStore()
	downloadCtrl.websocketStore = store

	// setup the websocket upgrader
	ws, err := provider.NewWebsocket(
		opt.Origins,
	)
	if err != nil {
		return nil, fmt.Errorf("provider.NewWebsocket: %v", err)
	}
	downloadCtrl.websocket = ws

	return downloadCtrl, nil
}

// Close closes the task client and the Pub/Sub helper
func (c *DownloadController) Close() error {
	if err := c.taskClient.Close(); err != nil {
		return fmt.Errorf("cloudtasks.Close: %v", err)
	}

	c.sub.Close()
	return nil
}
