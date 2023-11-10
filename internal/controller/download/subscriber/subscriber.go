// Package subscribe is a helper that interacts with Pub/Sub Cloud feature.
// It stores active websockets. When message is received, it notifies the
// concerned websocket if any.
package subscriber

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

// Subscriber interacts with Pub/Sub through a subscription. It holds the
// receive callback context and its cancel callback.
// It also has its own logger.
type Subscriber interface {

	// Close uses the cancel callback to stop receiving new messages from the
	// subscription.
	Close()

	// Listen uses sub.Receive to receive new messages from the subscription.
	// When a message is received, the configured callback is called.
	Listen() error

	// NewJobStatus converts a received Pub/Sub message into a jobBody.
	// It reads the message attributes to retrieve a code, a status and the ordering
	// key. It also parses the message data as JSON into a jobStatus entry.
	NewJobStatus(pMsg *pubsub.Message) (*JobStatus, error)
}

type subscriberImpl struct {
	// The Pub/Sub subscription used to receive messages
	sub *pubsub.Subscription

	// The context used to receive messages
	ctx context.Context

	// The callback to stop receiving messages
	cancel context.CancelFunc

	// The subscriber logger
	logger *log.Logger

	// The callback to use when a message is received
	onReceive func(ctx context.Context, message *pubsub.Message)
}

// SubscriberOptions holds the configuration to build a new subscriber
type SubscriberOptions struct {
	OnReceive      func(ctx context.Context, message *pubsub.Message)
	ProjectID      string
	SubscriptionID string
	LoggerBase     *log.Logger
}

// NewSubscriber builds a new Pub/Sub subscriber. It retrieve the subscription
// using the project ID and a subscription ID. The websocket store is
// initialized. The logger is created using the base logger configuration.
func NewSubscriber(opt SubscriberOptions) (*subscriberImpl, error) {

	client, err := pubsub.NewClient(context.Background(), opt.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	sub := client.Subscription(opt.SubscriptionID)

	ctx, cancel := context.WithCancel(context.Background())

	loggerPrefix := fmt.Sprintf("%s- [Pub/Sub] ", opt.LoggerBase.Prefix())
	logger := log.New(opt.LoggerBase.Writer(), loggerPrefix, opt.LoggerBase.Flags())

	return &subscriberImpl{
		sub:       sub,
		ctx:       ctx,
		cancel:    cancel,
		logger:    logger,
		onReceive: opt.OnReceive,
	}, nil
}

func (s *subscriberImpl) Close() {
	s.cancel()
	s.logger.Println("stopped listening for pubsub messages")
}

func (s *subscriberImpl) Listen() error {

	s.logger.Println("listening for pubsub messages..")
	err := s.sub.Receive(s.ctx, s.onReceive)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}
