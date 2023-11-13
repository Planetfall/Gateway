package mocks

import (
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/stretchr/testify/mock"
)

type ProviderMock struct {
	mock.Mock
}

func NewProviderMock(client task.Client) task.Provider {
	p := &ProviderMock{}
	p.On("NewClient").Return(client, nil)
	return p
}

func (m *ProviderMock) NewClient() (task.Client, error) {
	args := m.Called()
	return args.Get(0).(task.Client), args.Error(1)
}
