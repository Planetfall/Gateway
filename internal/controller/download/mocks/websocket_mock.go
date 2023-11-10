package mocks

import (
	"net"
	"net/http"

	"github.com/planetfall/gateway/internal/controller/download/websocket"
	"github.com/stretchr/testify/mock"
)

func NewWebsocketMock() websocket.Websocket {
	return &WebsocketMock{}
}

type WebsocketMock struct {
	mock.Mock
}

func (m *WebsocketMock) Upgrade(w http.ResponseWriter, r *http.Request,
	h http.Header) (websocket.Conn, error) {

	args := m.Called()
	return args.Get(0).(websocket.Conn), args.Error(1)
}

func (m *WebsocketMock) IsClosed(err error) bool {
	args := m.Called()
	return args.Bool(0)
}

func NewConnMock() websocket.Conn {
	return &ConnMock{}
}

type ConnMock struct {
	mock.Mock
}

func (m *ConnMock) ReadJSON(p interface{}) error {
	args := m.Called()
	return args.Error(0)
}

func (m *ConnMock) WriteJSON(p interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *ConnMock) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *ConnMock) Close() error {
	args := m.Called()
	return args.Error(0)
}
