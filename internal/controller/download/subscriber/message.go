package subscriber

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
)

// JobBody is the format of the job output retrieved from the Pub/Sub message
// Example:
// { "message": "job is still in progress", "progress": 35 }
type JobBody struct {
	Message  string `json:"message"`
	Progress int    `json:"progress"`
}

// JobStatus holds the jobStatus and the JobBody
type JobStatus struct {
	// The job output
	Body JobBody `json:"body"`

	// The code attribute of the Pub/Sub message
	Code int `json:"code"`

	// The status attribute of the Pub/Sub message
	Status string `json:"status"`

	// The ordering key of the Pub/Sub message
	OrderingKey string `json:"ordering_key"`
}

func (*subscriberImpl) NewJobStatus(
	pMsg *pubsub.Message) (*JobStatus, error) {

	// parse body
	var jBody JobBody
	if err := json.Unmarshal(pMsg.Data, &jBody); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	// parse code attribute
	codeStr, exists := pMsg.Attributes["code"]
	if !exists {
		return nil, fmt.Errorf("message.Attributes: 'code' not found")
	}

	// convert code to int
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return nil, fmt.Errorf("strconv.Atoi: %v", err)
	}

	// parse status attribute
	status, exists := pMsg.Attributes["status"]
	if !exists {
		return nil, fmt.Errorf("message.Attributes: 'status' not found")
	}

	// parse ordering key
	orderingKey := pMsg.OrderingKey

	return &JobStatus{
		Body:        jBody,
		Code:        code,
		Status:      status,
		OrderingKey: orderingKey,
	}, nil
}
