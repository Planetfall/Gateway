package controller

import (
	"context"
	"time"
)

const defaultTimeout = 10 * time.Second

// base controller type providing utils
type Controller struct {
	errorReportCallback func(err error)
}

type ControllerOptions struct {
	Host                string
	Audience            string
	Insecure            bool
	ErrorReportCallback func(err error)
}

func (c *Controller) getContext(
	timeout time.Duration) (context.Context, context.CancelFunc) {

	return context.WithTimeout(context.Background(), timeout)
}
