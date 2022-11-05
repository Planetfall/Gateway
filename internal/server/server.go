package server

import (
	"fmt"
	"log"

	"github.com/planetfall/gateway/internal/router"
)

type Server struct {
	port        string
	serviceName string
	projectID   string
	router      *router.Router
	gcloud      *Gcloud
}

type ServerOptions struct {
	ServiceName string
	ProjectID   string
	Port        string
	SvcConfig   router.ServicesConfig
	Insecure    bool
}

// @title          Gateway Front API
// @version        0.0.1
// @description    This application provides a front gateway allowing you
// @description    to interact with multiple GRPC microservices hosted
// @description    in Google Cloud
// @termsOfService No terms
// @contact.name   Support
// @contact.email  florian.charpentier67@gmail.com
// @license.name   MIT
// @license.url    http://opensource.org/licenses/MIT
// @host           api.dadard.fr
// @BasePath       /
// @accept         json
// @produce        json
// @schemes        https
func NewServer(opt ServerOptions) (*Server, error) {

	gcloud, err := NewGcloud(opt.ServiceName, opt.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("server.NewGcloud: %v", err)
	}

	router, err := router.NewRouter(
		router.RouterOptions{
			ProjectID:           opt.ProjectID,
			Insecure:            opt.Insecure,
			ErrorReportCallback: gcloud.ErrorReport,
			ConfigMap:           opt.SvcConfig,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("router.NewRouter: %v", err)
	}

	return &Server{
		port:        opt.Port,
		serviceName: opt.ServiceName,
		projectID:   opt.ProjectID,

		gcloud: gcloud,
		router: router,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("listening on port %s", s.port)
	if err := s.router.Run(":" + s.port); err != nil {
		return fmt.Errorf("router.Run: %v", err)
	}

	return nil
}

func (s *Server) Close() error {
	if err := s.gcloud.Close(); err != nil {
		return fmt.Errorf("gcloud.Close: %v", err)
	}

	if err := s.router.Close(); err != nil {
		return fmt.Errorf("router.Close: %v", err)
	}

	return nil
}
