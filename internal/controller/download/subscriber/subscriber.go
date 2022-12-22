package subscriber

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	client *pubsub.Client
	sub    *pubsub.Subscription

	listenContext context.Context
	listenCancel  context.CancelFunc

	wsStore websocketStore
}

func NewSubscriber(
	projectID string, subscriptionID string) (*Subscriber, error) {

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subscriptionID)

	lContext, lCancel := context.WithCancel(context.Background())

	return &Subscriber{
		sub:    sub,
		client: client,

		listenContext: lContext,
		listenCancel:  lCancel,

		wsStore: websocketStore(map[*websocket.Conn][]jobKey{}),
	}, nil
}

func (s *Subscriber) Close() {
	s.listenCancel()
	log.Println("[SUB] stopped listening for pubsub messages")
}

// wrapper around sub.Receive, notifies ws when received
func (s *Subscriber) Listen() error {

	log.Println("[SUB] listening for pubsub messages..")
	err := s.sub.Receive(
		s.listenContext,
		func(ctx context.Context, pMsg *pubsub.Message) {

			defer pMsg.Ack()

			// parse pMsg content
			jStatus, err := s.parsePubsubMessage(pMsg)
			if err != nil {
				log.Println(fmt.Errorf("[SUB] Subscriber.parsePubsubMessage: %v", err))
				return
			}

			log.Printf("[SUB] received msg %s - %d", jStatus.OrderingKey, jStatus.Code)

			// retrieve the websocket using the message ordering key
			ws, err := s.getWebsocketFromKey(pMsg.OrderingKey)
			if err != nil {
				log.Println(fmt.Errorf("[SUB] subscriber.getWebSocket: %v", err))
				return
			}

			// notify to ws
			if err := ws.WriteJSON(&jStatus); err != nil {
				log.Println(fmt.Errorf("[SUB] websocket.WriteMessage: %v", err))
			}
		})
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}
