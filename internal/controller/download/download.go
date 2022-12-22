package download

import (
	"context"
	"fmt"
	"log"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"

	"github.com/planetfall/gateway/internal/controller"
)

type DownloadController struct {
	controller.Controller

	jobUrl     string
	queuePath  string
	taskClient *cloudtasks.Client

	sub *subscriber.Subscriber
}

type DownloadControllerOptions struct {
	controller.ControllerOptions

	ProjectID string

	LocationID string
	QueueID    string

	SubscriptionID string
}

func NewDownloadController(opt DownloadControllerOptions) (*DownloadController, error) {

	ctx := context.Background()

	ctrl := controller.Controller{
		opt.ErrorReportCallback,
	}

	// task client setup
	jobUrl := fmt.Sprintf("%s/download/url", opt.Host)

	queuePath := fmt.Sprintf(
		"projects/%s/locations/%s/queues/%s",
		opt.ProjectID, opt.LocationID, opt.QueueID,
	)

	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.NewClient: %v", err)
	}

	// pubsub client setup
	sub, err := subscriber.NewSubscriber(opt.ProjectID, opt.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("newDownloadSub: %v", err)
	}
	// starts listening for pubsub messages
	go func() {
		if err := sub.Listen(); err != nil {
			log.Println(fmt.Errorf("[SUB] subscriber.Listen: %v", err))
		}
	}()

	return &DownloadController{
		ctrl,

		jobUrl,
		queuePath,
		taskClient,

		sub,
	}, nil
}

func (c *DownloadController) Close() error {
	if err := c.taskClient.Close(); err != nil {
		return fmt.Errorf("cloudtasks.Close: %v")
	}

	c.sub.Close()

	return nil
}
