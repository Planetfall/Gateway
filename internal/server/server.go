package server

import (
	"fmt"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port        string
	serviceName string
	router      *gin.Engine

	errorReporting *errorreporting.Client
	metadataClient *metadata.Client

	conns connections
	cls   *clients
}

func NewServer(
	env string, serviceName string, port string,
	connCfgList ConnectionConfigList,
) (*Server, error) {

	log.Println("setting up connections to services...")
	insecure := false
	if env == Development {
		// DEV env enforce the GRPC connections to be insecure
		insecure = true
		log.Println("INSECURE connections activated (no TLS, no token)")
		log.Println("DO NOT USE INSECURE IN PRODUCTION")
	}
	conns, err := newConnections(connCfgList, insecure)
	if err != nil {
		log.Println("failed setting up connections")
		return nil, err
	}
	cls := newClients(conns)

	var serv *Server
	switch env {
	case Development:
		serv, err = newServerDevelopment(serviceName, port, conns, cls)
		break
	case Production:
		serv, err = newServerProduction(serviceName, port, conns, cls)
		break
	default:
		return nil, fmt.Errorf(
			"failed to create server with unsupported env: %s", env)
	}

	if err != nil {
		log.Println("failed setting up the server")
		return nil, err
	}

	r := gin.Default()
	r.GET("/music-researcher/search", serv.musicSearch)

	serv.router = r
	return serv, nil
}

func (s *Server) Start() {
	log.Printf("listening on port %s", s.port)
	s.router.Run(":" + s.port)
}

func (s *Server) Close() error {
	// closing grpc connections
	if err := s.conns.Close(); err != nil {
		return err
	}

	// closing error reportings
	if s.errorReporting != nil {
		if err := s.errorReporting.Close(); err != nil {
			return err
		}
	}
	return nil
}
