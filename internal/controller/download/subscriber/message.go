package subscriber

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
)

// jobStatus is the format of the job output retrieved from the Pub/Sub message
type jobStatus struct {
	Message  string `json:"message"`
	Progress int    `json:"progress"`
}

// jobBody holds the jobStatus and the job metadatas
type jobBody struct {
	// The job output
	Body jobStatus `json:"body"`

	// The code attribute of the Pub/Sub message
	Code int `json:"code"`

	// The status attribute of the Pub/Sub message
	Status string `json:"status"`

	// The ordering key of the Pub/Sub message
	OrderingKey string `json:"ordering_key"`
}

// parsePubsubMessage converts a received Pub/Sub message into a jobBody.
// It reads the message attributes to retrieve a code, a status and the ordering
// key. It also parses the message data as JSON into a jobStatus entry.
func (*Subscriber) parsePubsubMessage(
	pMsg *pubsub.Message) (*jobBody, error) {

	// parse body
	var jBody jobStatus
	if err := json.Unmarshal(pMsg.Data, &jBody); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	// parse code attribute
	codeStr, exists := pMsg.Attributes["code"]
	if !exists {
		return nil, fmt.Errorf("pMsg.Attributes: 'code' not found")
	}

	// convert code to int
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return nil, fmt.Errorf("strconv.Atoi: %v", err)
	}

	// parse status attribute
	status, exists := pMsg.Attributes["status"]
	if !exists {
		return nil, fmt.Errorf("pMsg.Attributes: 'status' not found")
	}

	// parse ordering key
	orderingKey := pMsg.OrderingKey

	return &jobBody{
		Body:        jBody,
		Code:        code,
		Status:      status,
		OrderingKey: orderingKey,
	}, nil
}
