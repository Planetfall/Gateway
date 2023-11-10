package grpc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/planetfall/gateway/internal/connection/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	grpcG "google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

type tokenSourceMock struct {
	mock.Mock
}

func (m *tokenSourceMock) Token() (*oauth2.Token, error) {
	args := m.Called()
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func TestAuthenticateContext_withInsecureFalse(t *testing.T) {
	// given
	optGiven := grpc.ConnectionOptions{
		Target:   "target",
		Insecure: false,
	}

	c, err := grpc.NewConnection(optGiven)
	assert.Nil(t, err)

	// when
	ctxGiven := context.Background()
	ctxActual, err := c.AuthenticateContext(ctxGiven)
	assert.Nil(t, err)

	// then
	assert.NotEqual(t, ctxGiven, ctxActual)

	// check the context has the bearer token
	ctxActualData, found := grpcMetadata.FromOutgoingContext(ctxActual)
	assert.True(t, found)
	ctxActualAuthorization, found := ctxActualData["authorization"]
	assert.True(t, found)
	assert.Len(t, ctxActualAuthorization, 1)
	assert.Contains(t, ctxActualAuthorization[0], "Bearer")

	err = c.Close()
	assert.Nil(t, err)
}

func TestAuthenticateContext_withInsecureTrue(t *testing.T) {
	// given
	optGiven := grpc.ConnectionOptions{
		Target:   "target",
		Insecure: true,
	}

	c, err := grpc.NewConnection(optGiven)
	assert.Nil(t, err)

	// when
	ctxGiven := context.Background()
	ctxActual, err := c.AuthenticateContext(ctxGiven)
	assert.Nil(t, err)

	assert.Equal(t, ctxActual, ctxGiven)

	err = c.Close()
	assert.Nil(t, err)
}

func TestAuthenticateContext_withTokenSourceError(t *testing.T) {
	// given
	providerGiven := &providerMock{}
	tokenSourceGiven := &tokenSourceMock{}
	optGiven := grpc.ConnectionOptions{
		Target:   "target",
		Insecure: false,
		Provider: providerGiven,
	}

	providerGiven.
		On("NewClient", mock.Anything, mock.Anything).
		Return(&grpcG.ClientConn{}, nil)
	providerGiven.
		On("NewTokenSource", mock.Anything, mock.Anything).
		Return(tokenSourceGiven, nil)

	c, err := grpc.NewConnection(optGiven)
	assert.Nil(t, err)

	errMessageGiven := "test token error"
	tokenSourceGiven.
		On("Token").
		Return(&oauth2.Token{}, fmt.Errorf(errMessageGiven))

	// when
	ctxGiven := context.Background()
	_, err = c.AuthenticateContext(ctxGiven)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), errMessageGiven)
}
