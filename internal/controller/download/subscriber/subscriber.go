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
// sub.Receive context and its cancel callback. It also stores the active
// websockets, so that it can notify them when a message is received.
// It also has its own logger.
type Subscriber struct {
	// The Pub/Sub subscription used to receive messages
	sub *pubsub.Subscription

	// The context used to receive messages
	ctx context.Context

	// The callback to stop receiving messages
	cancel context.CancelFunc

	// The store which holds the active websockets and their job keys.
	// It is used to retrieve the concerned websocket when a message is
	// received.
	wsStore websocketStore

	// The subscriber logger
	logger *log.Logger
}

// NewSubscriber builds a new Pub/Sub subscriber. It retrieve the subscription
// using the project ID and a subscription ID. The websocket store is
// initialized. The logger is created using the base logger configuration.
func NewSubscriber(
	projectID string, subscriptionID string,
	baseLogger *log.Logger) (*Subscriber, error) {

	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	sub := client.Subscription(subscriptionID)

	ctx, cancel := context.WithCancel(context.Background())
	wsStore := make(websocketStore, 0)

	loggerPrefix := fmt.Sprintf("%s- [Pub/Sub] ", baseLogger.Prefix())
	logger := log.New(baseLogger.Writer(), loggerPrefix, baseLogger.Flags())

	return &Subscriber{
		sub: sub,

		ctx:     ctx,
		cancel:  cancel,
		wsStore: wsStore,

		logger: logger,
	}, nil
}

// Close uses the cancel callback to stop receiving new messages from the
// subscription.
func (s *Subscriber) Close() {
	s.cancel()
	s.logger.Println("stopped listening for pubsub messages")
}

// ReceiveCallback is called when a message is received from the subscription.
// It ensures that the message is acknowledged.
// The received message is parsed. Then, the calling websocket is retrieved in
// the store using the job ordering key. The parsed message is written on this
// websocket.
func (s *Subscriber) ReceiveCallback(
	ctx context.Context, message *pubsub.Message) {

	defer message.Ack()

	// parse pMsg content
	jobStatus, err := s.parsePubsubMessage(message)
	if err != nil {
		s.logger.Println(fmt.Errorf("Subscriber.parsePubsubMessage: %v", err))
		return
	}

	s.logger.Printf("received msg %s - %d", jobStatus.OrderingKey, jobStatus.Code)

	// retrieve the websocket using the message ordering key
	ws, err := s.getWebsocketFromKey(message.OrderingKey)
	if err != nil {
		s.logger.Println(fmt.Errorf("subscriber.getWebSocket: %v", err))
		return
	}

	// notify to ws
	if err := ws.WriteJSON(&jobStatus); err != nil {
		s.logger.Println(fmt.Errorf("websocket.WriteMessage: %v", err))
	}
}

// Listen uses sub.Receive to receive new messages from the subscription.
// When a message is received, subscriber.ReceiveCallback is called.
func (s *Subscriber) Listen() error {

	s.logger.Println("listening for pubsub messages..")
	err := s.sub.Receive(s.ctx, s.ReceiveCallback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}
