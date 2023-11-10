// Package controller provides a base helper type.
// It holds some common behavior for all controllers to use.
package controller

import (
	"context"
	"log"
	"time"
)

// The defaut context timeout
const DefaultTimeout = 10 * time.Second

// Controller is a base type for actual controllers.
// It is a common base and holds some behavior like error handling and
// context management
type Controller struct {
	// name of the controller
	name string

	// Target indicates on which URL or host the controller should connect
	Target string

	// ReportError is a callback, that is used for error handling at the
	// Service level.
	ReportError func(err error)

	// The logger to use for the controller. Having a logger per controller
	// allows to split and easily filters the console output.
	Logger *log.Logger
}

// ControllerOptions holds the base parameters needed to build a Controller
type ControllerOptions struct {
	// Name builder parameter
	Name string

	// Target builder parameter
	Target string

	// ReportError builder parameter
	ReportError func(err error)

	// Logger builder parameter
	Logger *log.Logger
}

// NewController builds a new controller. It setup a new logger using the
// provided prefix uppercased.
func NewController(opt ControllerOptions) Controller {

	return Controller{
		name:        opt.Name,
		Target:      opt.Target,
		ReportError: opt.ReportError,
		Logger:      opt.Logger,
	}
}

func (c *Controller) Name() string {
	return c.name
}

// GetContext is a wrapper around context.WithTimeout, with a default timeout
// value. An optional custom timeout value can be provided. Only the first
// custom timeout value will be used. If more than one parameter is given, it
// will fallback to the defaultTimeout value.
func (c *Controller) GetContext(
	timeout ...time.Duration) (context.Context, context.CancelFunc) {

	t := DefaultTimeout
	if len(timeout) == 1 {
		t = timeout[0]
	}
	return context.WithTimeout(context.Background(), t)
}
