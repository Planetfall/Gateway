package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorMessage is the standard error response body sent back by a Controller
// in case of error. This standard format can be useful for frontend
// application in order to report back this kind of error to the user.
type errorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// logAndReport process the provided error with some steps:
//
//  1. Logs the error in the console
//  2. Sends it as a HTTP response using the errorMessage format
//  3. Use the controller ReportError callback
//
// This is a common workflow used by all controllers when encountering an error.
// It avoid exposing sensitive error message in the HTTP response.
func (c *Controller) logAndReport(
	err error, g *gin.Context, status int, message string) {

	// log the error in the console
	c.Logger.Println(err)

	// respond as HTTP with a general message
	g.JSON(status, errorMessage{
		Status:  status,
		Message: message,
	})

	// report the error to the server
	c.ReportError(err)
}

// BadRequest uses logAndReport with a http.StatusBadRequest status and a proper
// bad request message
func (c *Controller) BadRequest(err error, g *gin.Context) {
	c.logAndReport(err, g,
		http.StatusBadRequest, "Wrong parameters supplied")
}

// InternalError uses logAndReport with a http.StatusInternalServerError and a
// proper internal error message
func (c *Controller) InternalError(err error, g *gin.Context) {
	c.logAndReport(err, g,
		http.StatusInternalServerError, "Something went wrong on my side")
}
