package subscriber_test

import (
	"strconv"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestNewJobStatus(t *testing.T) {

	s := getSubscriber(t)

	givenCode := 200
	givenStatus := "OK"
	givenKey := "key"
	givenMessage := &pubsub.Message{
		Data: []byte(`{
			"message": "job is in progress",
			"progress": 32
		}`),
		Attributes: map[string]string{
			"code":   strconv.Itoa(givenCode),
			"status": givenStatus,
		},
		OrderingKey: givenKey,
	}

	actualJobStatus, err := s.NewJobStatus(givenMessage)
	assert.Nil(t, err)

	assert.Equal(t, givenCode, actualJobStatus.Code)
	assert.Equal(t, givenStatus, actualJobStatus.Status)
	assert.Equal(t, givenKey, actualJobStatus.OrderingKey)

	assert.Equal(t, "job is in progress", actualJobStatus.Body.Message)
	assert.Equal(t, 32, actualJobStatus.Body.Progress)
}

func TestNewJobStatus_withInvalidJson(t *testing.T) {

	s := getSubscriber(t)
	givenMessage := &pubsub.Message{
		Data: []byte(`invalid JSON`),
	}

	_, err := s.NewJobStatus(givenMessage)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json.Unmarshal")
}

func TestNewJobsStatus_withNoCode(t *testing.T) {

	s := getSubscriber(t)
	givenMessage := &pubsub.Message{
		Data: []byte(`{
			"message": "job in progress",
			"progress": 23
		}`),
	}

	_, err := s.NewJobStatus(givenMessage)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "message.Attributes: 'code' not found")
}

func TestNewJobsStatus_withInvalidCode(t *testing.T) {
	s := getSubscriber(t)
	givenMessage := &pubsub.Message{
		Data: []byte(`{
			"message": "job in progress",
			"progress": 23
		}`),
		Attributes: map[string]string{
			"code": "invalid code",
		},
	}

	_, err := s.NewJobStatus(givenMessage)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "strconv.Atoi")
}

func TestNewJobStatus_withNoStatus(t *testing.T) {
	s := getSubscriber(t)
	givenMessage := &pubsub.Message{
		Data: []byte(`{
			"message": "job in progress",
			"progress": 23
		}`),
		Attributes: map[string]string{
			"code": "200",
		},
	}

	_, err := s.NewJobStatus(givenMessage)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "message.Attributes: 'status' not found")
}
