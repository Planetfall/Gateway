package websocket

import (
	"fmt"
	"net"
)

// Conn is the type which represent an connection with the client.
// grpc.Conn is an example of implementation.
type Conn interface {
	// ReadJson reads the current input.
	ReadJSON(p interface{}) error

	// WriteJson writes back to the client.
	WriteJSON(p interface{}) error

	// RemoteAddr returns the remote address of the client.
	RemoteAddr() net.Addr

	// Close closes the connection with the client.
	Close() error
}

type Status string

var StatusOK Status = "ok"
var StatusError Status = "error"

type StatusBody struct {
	Status  Status      `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

func WriteStatus(conn Conn, status Status, message string, body interface{}) error {
	payload := StatusBody{
		Status:  status,
		Message: message,
		Body:    body,
	}

	if err := conn.WriteJSON(&payload); err != nil {
		return fmt.Errorf("connection.WriteJSON: %v", err)
	}

	return nil
}
