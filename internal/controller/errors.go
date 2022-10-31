package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (c *Controller) formatError(
	err error, g *gin.Context, status int, message string) {

	log.Println(err)
	msg := errorMessage{
		Status:  status,
		Message: message,
	}
	g.JSON(status, msg)

	// use the callback provided by the Server
	c.errorReportCallback(err)
}

func (c *Controller) badRequest(err error, g *gin.Context) {
	c.formatError(err, g,
		http.StatusBadRequest, "Wrong parameters supplied")
}

func (c *Controller) internalError(err error, g *gin.Context) {
	c.formatError(err, g,
		http.StatusInternalServerError, "Something went wrong on my side")
}
