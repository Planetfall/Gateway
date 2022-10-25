package server

import (
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
	"github.com/gin-gonic/gin"
)

// errors
type errorMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (s *Server) formatError(err error, c *gin.Context, status int, message string) {
	log.Println(err)
	msg := errorMessage{
		Status:  status,
		Message: message,
	}
	c.JSON(status, msg)
	s.errorReporting.Report(errorreporting.Entry{
		Error: err,
	})
}

func (s *Server) badRequest(err error, c *gin.Context) {
	s.formatError(err, c,
		http.StatusBadRequest, "Wrong parameters supplied")
}

func (s *Server) internalError(err error, c *gin.Context) {
	s.formatError(err, c,
		http.StatusInternalServerError, "Something went wrong on my side")
}
