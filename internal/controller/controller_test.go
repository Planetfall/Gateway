package controller_test

import (
	"context"
	"log"
	"math"
	"testing"
	"time"

	"github.com/planetfall/gateway/internal/controller"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	// given
	nameGiven := "controller-name"
	reportGiven := func(err error) { t.Logf(err.Error()) }
	loggerGiven := log.Default()
	optGiven := controller.ControllerOptions{
		Name:        nameGiven,
		Target:      "target",
		ReportError: reportGiven,
		Logger:      loggerGiven,
	}

	// when
	c := controller.NewController(optGiven)

	// then
	assert.Equal(t, nameGiven, c.Name())
}

func TestGetContext(t *testing.T) {
	// given
	c := &controller.Controller{}

	// when
	ctx, cancel := c.GetContext()
	cancel()

	// then
	assert.Equal(t, ctx.Err(), context.Canceled)
}

func TestGetContextAndWait(t *testing.T) {
	// given
	c := &controller.Controller{}

	// when
	ctx, _ := c.GetContext()

	// then
	waitTimeSeconds := math.Floor(controller.DefaultTimeout.Seconds()/2) + 1
	waitTime := time.Duration(waitTimeSeconds) * time.Second
	time.Sleep(waitTime)
	assert.Nil(t, ctx.Err())

	time.Sleep(waitTime)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}

func TestGetContextAndWait_withCustomTimeout(t *testing.T) {
	// given
	timeoutGiven := 6 * time.Second
	c := &controller.Controller{}

	// when
	ctx, _ := c.GetContext(timeoutGiven)

	// then
	waitTimeSeconds := math.Floor(timeoutGiven.Seconds()/2) + 1
	waitTime := time.Duration(waitTimeSeconds) * time.Second
	time.Sleep(waitTime)
	assert.Nil(t, ctx.Err())

	time.Sleep(waitTime)
	assert.Equal(t, context.DeadlineExceeded, ctx.Err())
}
