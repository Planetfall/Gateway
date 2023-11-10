package websocket

import "net"

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
