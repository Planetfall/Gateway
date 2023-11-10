package grpc_test

import (
	"fmt"
	"testing"

	"github.com/planetfall/gateway/internal/connection/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	grpcG "google.golang.org/grpc"
)

type providerMock struct {
	mock.Mock
}

func (m *providerMock) NewTokenSource(
	audience string, insecure bool) (grpc.TokenSource, error) {

	args := m.Called(audience, insecure)
	return args.Get(0).(grpc.TokenSource), args.Error(1)
}

func (m *providerMock) NewClient(
	target string, insecure bool) (*grpcG.ClientConn, error) {

	args := m.Called(target, insecure)
	return args.Get(0).(*grpcG.ClientConn), args.Error(1)
}

func TestNewConnection_withInsecureTrue(t *testing.T) {
	// given
	optGiven := grpc.ConnectionOptions{
		Target:   "target",
		Insecure: true,
	}

	// when
	c, err := grpc.NewConnection(optGiven)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.Client())

	err = c.Close()
	assert.Nil(t, err)
}

func TestNewConnection_withInsecureFalse(t *testing.T) {
	// given
	optGiven := grpc.ConnectionOptions{
		Target:   "target",
		Insecure: false,
	}

	// when
	c, err := grpc.NewConnection(optGiven)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, c)

	err = c.Close()
	assert.Nil(t, err)
}

func TestNewConnect_withClientError(t *testing.T) {
	// given
	providerGiven := &providerMock{}
	targetGiven := "target"
	insecureGiven := false
	optGiven := grpc.ConnectionOptions{
		Target:   targetGiven,
		Insecure: insecureGiven,
		Provider: providerGiven,
	}

	// when
	errMessageGiven := "test client error"
	providerGiven.
		On("NewClient", targetGiven, insecureGiven).
		Return(&grpcG.ClientConn{}, fmt.Errorf(errMessageGiven))

	// then
	c, err := grpc.NewConnection(optGiven)

	providerGiven.AssertExpectations(t)
	providerGiven.AssertNotCalled(t, "NewTokenSource")

	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMessageGiven)
}

func TestNewConnect_withTokenSourceError(t *testing.T) {
	// given
	providerGiven := &providerMock{}
	targetGiven := "target"
	insecureGiven := false
	optGiven := grpc.ConnectionOptions{
		Target:   targetGiven,
		Insecure: insecureGiven,
		Provider: providerGiven,
	}

	tokenSourceGiven := &tokenSourceMock{}

	// when
	providerGiven.
		On("NewClient", targetGiven, insecureGiven).
		Return(&grpcG.ClientConn{}, nil)
	errMessageGiven := "test token source error"
	providerGiven.
		On("NewTokenSource", mock.Anything, insecureGiven).
		Return(tokenSourceGiven, fmt.Errorf(errMessageGiven))

	// then
	c, err := grpc.NewConnection(optGiven)

	providerGiven.AssertExpectations(t)
	providerGiven.AssertNotCalled(t, "NewTokenSource")

	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMessageGiven)
}
