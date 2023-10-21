package controller_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/controller"
	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	// given
	writerGiven := httptest.NewRecorder()
	ginContextGiven, _ := gin.CreateTestContext(writerGiven)

	errorGiven := fmt.Errorf("test error")
	loggerGiven := log.Default()

	reportedGiven := false
	reportErrorGiven := func(err error) {
		t.Logf(err.Error())
		reportedGiven = true
	}

	c := &controller.Controller{
		Logger:      loggerGiven,
		ReportError: reportErrorGiven,
	}

	// when
	c.BadRequest(errorGiven, ginContextGiven)

	// then
	assert.Equal(t,
		http.StatusBadRequest, writerGiven.Result().StatusCode)
	assert.True(t, reportedGiven)
}

func TestInternalError(t *testing.T) {
	// given
	writerGiven := httptest.NewRecorder()
	ginContextGiven, _ := gin.CreateTestContext(writerGiven)

	errorGiven := fmt.Errorf("test error")
	loggerGiven := log.Default()

	reportedGiven := false
	reportErrorGiven := func(err error) {
		t.Logf(err.Error())
		reportedGiven = true
	}

	c := &controller.Controller{
		Logger:      loggerGiven,
		ReportError: reportErrorGiven,
	}

	// when
	c.InternalError(errorGiven, ginContextGiven)

	// then
	assert.Equal(t,
		http.StatusInternalServerError, writerGiven.Result().StatusCode)
	assert.True(t, reportedGiven)
}
