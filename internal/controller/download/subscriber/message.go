package subscriber

import (
	"encoding/json"
	"fmt"
	"strconv"

	"cloud.google.com/go/pubsub"
)

type jobBody struct {
	Message  string `json:"message"`
	Progress int    `json:"progress"`
}

type jobStatus struct {
	Body jobBody `json:"body"`

	Code        int    `json:"code"`
	Status      string `json:"status"`
	OrderingKey string `json:"ordering_key"`
}

func (*Subscriber) parsePubsubMessage(
	pMsg *pubsub.Message) (*jobStatus, error) {

	// parse body
	var jBody jobBody
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

	return &jobStatus{
		Body:        jBody,
		Code:        code,
		Status:      status,
		OrderingKey: orderingKey,
	}, nil
}
