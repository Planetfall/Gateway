package subscriber_test

import (
	"log"
	"testing"

	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/stretchr/testify/assert"
)

func getSubscriber(t *testing.T) subscriber.Subscriber {

	opt := subscriber.SubscriberOptions{
		ProjectID:  "project-id",
		LoggerBase: log.Default(),
	}
	s, err := subscriber.NewSubscriber(opt)
	assert.Nil(t, err)

	return s
}

func TestListen_withError(t *testing.T) {

	s := getSubscriber(t)

	err := s.Listen()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "sub.Receive")
}

func TestClose(t *testing.T) {

	s := getSubscriber(t)
	s.Close()

	err := s.Listen()
	assert.Nil(t, err)
}

func TestNewSubscriber_withNoProjectID(t *testing.T) {
	_, err := subscriber.NewSubscriber(subscriber.SubscriberOptions{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "pubsub.NewClient")
}
