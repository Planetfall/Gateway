package runserver

import (
	"fmt"

	"github.com/planetfall/gateway/internal/router"
	"github.com/planetfall/gateway/internal/server"
	"github.com/planetfall/gateway/pkg/config"
)

type RunServerOptions struct {
	ServiceName string
	ProjectID   string
	Port        string
	SvcConfig   config.ServiceConfigMap
	Insecure    bool

	// unused
	Environment string
}

func RunServer(opt RunServerOptions) error {
	services := router.ServicesConfig{}
	for service, conf := range opt.SvcConfig {
		services[service] = router.ServiceConfig(conf)
	}

	s, err := server.NewServer(
		server.ServerOptions{
			ServiceName: opt.ServiceName,
			ProjectID:   opt.ProjectID,
			Port:        opt.Port,
			SvcConfig:   services,
			Insecure:    opt.Insecure,
		},
	)
	if err != nil {
		return fmt.Errorf("server.NewServer: %v", err)
	}

	if err := s.Start(); err != nil {
		return fmt.Errorf("server.Start: %v", err)
	}

	return nil
}
