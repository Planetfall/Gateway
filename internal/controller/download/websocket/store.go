package websocket

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Key is a key representing a job for a connection.
type Key string

// Store stores the map between connections and its job keys.
type Store map[Conn][]Key

// NewStore builds a new store.
func NewStore() Store {
	return Store{}
}

// Generate a new job key from a websocket, using the client addr and
// the current timestamp
func (w Store) newWebsocketKey(ws Conn) Key {
	addr := ws.RemoteAddr().String()
	now := fmt.Sprintf("%v", time.Now())

	id := fmt.Sprintf("%s-%s", addr, now)
	sum := sha256.Sum256([]byte(id))
	key := fmt.Sprintf("%x", sum)

	return Key(key[:8])
}

// GetWebsocket provides the websocket associated with the provided job
// key
func (s Store) GetWebsocket(key Key) (Conn, error) {

	// loop over ws
	for ws, wsKeys := range s {

		// for each ws, check if it contains the filtered ordering key
		for _, wsKey := range wsKeys {
			if wsKey == key {
				return ws, nil
			}
		}
	}

	return nil, fmt.Errorf("key %s not found", key)
}

// RegisterWebsocket adds a new websocket in the store.
// It initialize its job key slice.
func (s Store) Register(ws Conn) error {

	if _, exists := s[ws]; exists {
		return fmt.Errorf("websocket already registered")
	}

	// initialize with empty array
	s[ws] = make([]Key, 0)

	return nil
}

// AddNewJob adds a new job key to a registered websocket.
func (s Store) AddNewJob(ws Conn) (Key, error) {

	if _, exists := s[ws]; !exists {
		return "", fmt.Errorf("websocket not registered")
	}

	newKey := s.newWebsocketKey(ws)
	s[ws] = append(s[ws], newKey)

	return newKey, nil
}

// UnregisterWebsocket removes a registered websocket from the store and removes
// all its job keys
func (s Store) Unregister(ws Conn) error {

	if _, exists := s[ws]; !exists {
		return fmt.Errorf("websocket not registered")
	}

	delete(s, ws)

	return nil
}
