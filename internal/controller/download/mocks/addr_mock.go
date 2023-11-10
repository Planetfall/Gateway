package mocks

import (
	"net"

	"github.com/stretchr/testify/mock"
)

type AddrMock struct {
	mock.Mock
}

func NewAddrMock(addr string) net.Addr {
	m := &AddrMock{}
	m.On("String").Return(addr)
	return m
}

func (m *AddrMock) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *AddrMock) Network() string {
	args := m.Called()
	return args.String(0)
}
