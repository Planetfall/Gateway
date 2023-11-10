package websocket

import (
	"fmt"
	"net/http"

	ws "github.com/gorilla/websocket"
)

type WebsocketOptions struct {
	Origins []string
}

type websocketImpl struct {
	upgrader ws.Upgrader
}

func NewWebsocket(opt WebsocketOptions) (Websocket, error) {

	checkOrigins := func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// fixme
		// no origin to verify (aka, not called by javascript)
		if origin == "" {
			return true
		}

		for _, authorizedOrigin := range opt.Origins {
			if authorizedOrigin == origin {
				return true
			}
		}
		return false
	}

	upgrader := ws.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin:     checkOrigins,
	}

	return &websocketImpl{
		upgrader: upgrader,
	}, nil
}

func (w *websocketImpl) Upgrade(writer http.ResponseWriter, r *http.Request,
	h http.Header) (Conn, error) {

	conn, err := w.upgrader.Upgrade(writer, r, h)
	if err != nil {
		return nil, fmt.Errorf("upgrader.Upgrade: %v", err)
	}

	return conn, nil
}

func (w *websocketImpl) IsClosed(err error) bool {

	return ws.IsCloseError(err,
		ws.CloseNormalClosure,
		ws.CloseAbnormalClosure)
}
