// Package service manage the service controllers and the HTTP server.
package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/planetfall/framework/pkg/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Service holds the server helper from the framework, the HTTP server/router,
// the list of active controllers and a logger.
//
// It manages the initialization of the controllers from the configuration, and
// closing them when it shuts down.
type Service struct {
	srv *server.Server
	g   *gin.Engine

	ctrlList []svcController
}

// swaggerUIRoute is the path to get the swagger UI.
const swaggerUIRoute = "/swagger-ui/*any"

// ServiceOptions holds the service builder parameters
type ServiceOptions struct {
	// Srv builer parameter
	Srv *server.Server

	// ProjectID builder parameter
	ProjectID string

	// Insecure builder parameter
	Insecure bool
}

//	@title			Gateway Front API
//	@version		0.0.1
//	@description	This application provides a front gateway allowing you
//	@description	to interact with multiple GRPC microservices hosted
//	@description	in Google Cloud
//	@termsOfService	No terms
//	@contact.name	Support
//	@contact.email	florian.charpentier67@gmail.com
//	@license.name	MIT
//	@license.url	http://opensource.org/licenses/MIT
//	@host			api.dadard.fr
//	@BasePath		/
//	@accept			json
//	@produce		json
//	@schemes		https
func NewService(opt ServiceOptions) (*Service, error) {
	g := gin.Default()

	// cors middleware
	gConfig := cors.DefaultConfig()
	gConfig.AllowAllOrigins = true
	g.Use(cors.New(gConfig))

	// initialize the service
	svc := &Service{
		srv:      opt.Srv,
		g:        g,
		ctrlList: make([]svcController, 0),
	}

	// build controllers
	for key, builder := range controllerBuilders {

		svc.srv.Logger.Printf("setup controller %s", key)

		keyUpper := strings.ToUpper(key)
		loggerPrefix := fmt.Sprintf(
			"%s[%s] ", opt.Srv.Logger.Prefix(), keyUpper)
		logger := log.New(
			os.Stdout, loggerPrefix, log.Ldate|log.Ltime)

		route := fmt.Sprintf("/%s", key)
		group := g.Group(route)

		opt := svcControllerOptions{
			cfgKey:              key,
			group:               group,
			reportErrorCallback: svc.reportErrorCallback,
			logger:              logger,
			insecure:            opt.Insecure,
			projectID:           opt.ProjectID,
		}

		ctrl, err := builder(opt)
		if err != nil {
			return nil, fmt.Errorf(
				"builder for controller %s failed: %v", key, err)
		}

		svc.ctrlList = append(svc.ctrlList, ctrl)
	}

	// setup the swagger route
	svc.g.GET(swaggerUIRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	return svc, nil
}

// Start runs the HTTP server.
// It runs until the Interupt signal are received. Then, the HTTP server is
// gracefully shutdown. The service.Close call is deferred.
func (s *Service) Start(port string) error {
	addr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.g,
	}

	defer s.close()

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.srv.Logger.Fatalf("ListenAndServe: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	s.srv.Logger.Println("stopping")

	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("srv.Shutdown: %v", err)
	}

	s.srv.Logger.Println("stopped")
	return nil
}

// close frees all resources from the service.
// It closes all the active controllers stored in ctrlList.
//
// If one controller fails the close, a warning message is printed. The flow
// is not interrupted, and the service tries to close the other controllers
// anyway.
func (s *Service) close() {
	for _, ctrl := range s.ctrlList {
		err := ctrl.Close()
		if err != nil {
			s.srv.Logger.Printf(
				"failed to close controller %s: %v", ctrl.Name(), err)
		} else {
			s.srv.Logger.Printf(
				"closed controller %s", ctrl.Name(),
			)
		}
	}
}
