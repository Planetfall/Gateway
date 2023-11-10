package search_test

import (
	"fmt"
	"testing"

	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/search"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNewSearchController(t *testing.T) {
	// given
	optGiven := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:   "name",
			Target: "target",
		},
		Insecure: true,
	}

	// when
	c, err := search.NewSearchController(optGiven)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestNewSearchController_withEmptyTarget_shouldFail(t *testing.T) {
	// given
	optGiven := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name: "name",
		},
		Insecure: true,
	}

	// when
	c, err := search.NewSearchController(optGiven)

	// then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "NewConnection")
	assert.Nil(t, c)
}

func TestClose(t *testing.T) {
	// given
	optGiven := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:   "name",
			Target: "target",
		},
		Insecure: true,
	}
	c, _ := search.NewSearchController(optGiven)

	// when
	err := c.Close()

	// then
	assert.Nil(t, err)
}

func TestClose_shouldFail(t *testing.T) {
	// given
	connGiven := &connectionMock{}
	optGiven := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:   "name",
			Target: "target",
		},
		Insecure: true,
		Conn:     connGiven,
	}

	errMessageGiven := "test close error"
	connGiven.On("Client").Return(&grpc.ClientConn{})
	connGiven.On("Close").Return(fmt.Errorf(errMessageGiven))
	c, _ := search.NewSearchController(optGiven)

	// when
	errActual := c.Close()

	// then
	connGiven.AssertExpectations(t)
	assert.NotNil(t, errActual)
	assert.Contains(t, errActual.Error(), errMessageGiven)
}
