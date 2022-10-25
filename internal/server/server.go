package server

import (
	"context"
	"log"
	"net/http"

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

	conns *connections
	cls   *clients
}

func NewServer(serviceName string, port string) (*Server, error) {
	ctx := context.Background()

	log.Println("initializing metadata client...")
	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		log.Printf("metadata.ProjectID: %v\n", err)
		return nil, err
	}

	log.Println("initializing error reporting...")
	errorReporting, err := errorreporting.NewClient(ctx, projectID,
		errorreporting.Config{
			ServiceName: serviceName,
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			}},
	)
	if err != nil {
		log.Printf("errorreporting.NewClient: %v\n", err)
		return nil, err
	}

	conns, err := newConnections()
	if err != nil {
		log.Printf("failed settings up connections: %v\n", err)
		return nil, err
	}
	cls := newClients(conns)

	serv := &Server{
		port:           port,
		router:         nil,
		errorReporting: errorReporting,
		metadataClient: metadataClient,
		serviceName:    serviceName,
		conns:          conns,
		cls:            cls,
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
