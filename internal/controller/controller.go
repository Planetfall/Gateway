package controller

import (
	"context"
	"time"
)

const DefaultTimeout = 10 * time.Second

// base controller type providing utils
type Controller struct {
	ErrorReportCallback func(err error)
}

type ControllerOptions struct {
	Host                string
	Audience            string
	Insecure            bool
	ErrorReportCallback func(err error)
}

func (c *Controller) GetContext(
	timeout time.Duration) (context.Context, context.CancelFunc) {

	return context.WithTimeout(context.Background(), timeout)
}
