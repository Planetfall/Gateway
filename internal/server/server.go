package server

import (
	"fmt"
	"log"

	"github.com/planetfall/gateway/internal/router"
	"github.com/planetfall/gateway/pkg/options"
)

type Server struct {
	port        string
	serviceName string
	router      *router.Router
	gcloud      *Gcloud
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
func NewServer(opt options.ServerOptions) (*Server, error) {

	gcloud, err := NewGcloud(opt.Env, opt.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("server.NewGcloud: %v", err)
	}

	router, err := router.NewRouter(
		opt.SvcConfigMap, gcloud.ErrorReport, opt.Insecure,
	)
	if err != nil {
		return nil, fmt.Errorf("router.NewRouter: %v", err)
	}

	return &Server{
		port:        opt.Port,
		serviceName: opt.ServiceName,
		router:      router,
		gcloud:      gcloud,
	}, nil
}

func (s *Server) Start() {
	log.Printf("listening on port %s", s.port)
	s.router.Run(":" + s.port)
}

func (s *Server) Close() error {
	// fixme
	return nil
}
