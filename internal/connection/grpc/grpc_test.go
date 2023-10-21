package grpc_test

import (
	"testing"

	"github.com/planetfall/gateway/internal/connection/grpc"
	"github.com/stretchr/testify/assert"
)

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
	assert.NotNil(t, c.GrpcConn())

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
