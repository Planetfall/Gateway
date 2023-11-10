package websocket

import (
	"net/http"
)

// Websocket is responsible for creating new websocket connections.
type Websocket interface {

	// Upgrade creates a new connection from the current HTTP context.
	Upgrade(w http.ResponseWriter, r *http.Request,
		h http.Header) (Conn, error)

	// IsClosed checks if the given error is related to a connection closure.
	IsClosed(err error) bool
}
