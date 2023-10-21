package subscriber

import (
	"fmt"
	"time"

	"crypto/sha256"

	"github.com/gorilla/websocket"
)

// jobKey holds the key of the job started in a websocket
type jobKey string

// websocketStore holds the current active websockets. For each, it also stores
// the current job keys started in those websockets.
type websocketStore map[*websocket.Conn][]jobKey

// Generate a new job key from a websocket, using the client addr and
// the current timestamp
func (s *Subscriber) newJobKey(ws *websocket.Conn) jobKey {
	addr := ws.RemoteAddr().String()
	now := fmt.Sprintf("%v", time.Now())

	id := fmt.Sprintf("%s-%s", addr, now)
	sum := sha256.Sum256([]byte(id))
	key := fmt.Sprintf("%x", sum)

	return jobKey(key[:8])
}

// getWebsocketFromKey provides the websocket associated with the provided job
// key
func (s *Subscriber) getWebsocketFromKey(jobKeyStr string) (*websocket.Conn, error) {

	jKeyFilter := jobKey(jobKeyStr)
	// loop over ws
	for ws, jKeys := range s.wsStore {

		// for each ws, check if it contains the filtered ordering key
		for _, jKey := range jKeys {
			if jKey == jKeyFilter {
				return ws, nil
			}
		}
	}

	return nil, fmt.Errorf("job key %s not found", jobKeyStr)
}

// RegisterWebsocket adds a new websocket in the store.
// It initialize its job key slice.
func (s *Subscriber) RegisterWebsocket(ws *websocket.Conn) error {

	if _, exists := s.wsStore[ws]; exists {
		return fmt.Errorf("websocket already registered")
	}

	// initialize with empty array
	s.wsStore[ws] = make([]jobKey, 0)

	return nil
}

// AddNewJob adds a new job key to a registered websocket.
func (s *Subscriber) AddNewJob(ws *websocket.Conn) (jobKey, error) {

	if _, exists := s.wsStore[ws]; !exists {
		return "", fmt.Errorf("websocket not registered")
	}

	newJKey := s.newJobKey(ws)
	s.wsStore[ws] = append(s.wsStore[ws], newJKey)

	return newJKey, nil
}

// UnregisterWebsocket removes a registered websocket from the store and removes
// all its job keys
func (s *Subscriber) UnregisterWebsocket(ws *websocket.Conn) error {

	if _, exists := s.wsStore[ws]; !exists {
		return fmt.Errorf("websocket not registered")
	}

	delete(s.wsStore, ws)

	return nil
}
