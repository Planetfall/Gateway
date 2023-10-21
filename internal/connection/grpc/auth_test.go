package grpc_test

import (
	"context"
	"testing"

	"github.com/planetfall/gateway/internal/connection/grpc"
	"github.com/stretchr/testify/assert"
	grpcMetadata "google.golang.org/grpc/metadata"
)

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
