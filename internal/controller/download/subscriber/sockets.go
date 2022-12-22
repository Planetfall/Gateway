package subscriber

import (
	"fmt"
	"log"
	"time"

	"crypto/sha256"

	"github.com/gorilla/websocket"
)

type jobKey string
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

// Returns the websocket which has the jKeyStr associate
func (s *Subscriber) getWebsocketFromKey(jKeyStr string) (*websocket.Conn, error) {

	jKeyFilter := jobKey(jKeyStr)
	// loop over ws
	for ws, jKeys := range s.wsStore {

		// for each ws, check if it contains the filtered ordering key
		for _, jKey := range jKeys {
			if jKey == jKeyFilter {
				return ws, nil
			}
		}
	}

	return nil, fmt.Errorf("key %s not found", jKeyStr)
}

// Add a new websocket to the store and initialize its associate
// job keys slice
func (s *Subscriber) RegisterWebsocket(ws *websocket.Conn) error {

	if _, exists := s.wsStore[ws]; exists {
		return fmt.Errorf("websocket already registered")
	}

	// initialize with empty array
	s.wsStore[ws] = make([]jobKey, 0)

	log.Println(s.wsStore)

	return nil
}

// Add a new job key to an already registered websocket
func (s *Subscriber) AddNewJob(ws *websocket.Conn) (jobKey, error) {

	if _, exists := s.wsStore[ws]; !exists {
		return "", fmt.Errorf("websocket not registered")
	}

	newJKey := s.newJobKey(ws)
	s.wsStore[ws] = append(s.wsStore[ws], newJKey)

	return newJKey, nil
}

// Remove a registered websocket from the store and remove all
// associate job keys
func (s *Subscriber) UnregisterWebsocket(ws *websocket.Conn) error {

	if _, exists := s.wsStore[ws]; !exists {
		return fmt.Errorf("websocket not registered")
	}

	delete(s.wsStore, ws)
	log.Println(s.wsStore)

	return nil
}
