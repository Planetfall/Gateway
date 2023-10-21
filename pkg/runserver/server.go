package runserver

import (
	"fmt"

	"github.com/planetfall/framework/pkg/config"
	"github.com/planetfall/framework/pkg/server"
	"github.com/spf13/viper"
)

func getServer(cfg *config.Config) (*server.Server, error) {

	serviceName := viper.GetString(serviceFlag)
	s, err := server.NewServer(cfg, serviceName)
	if err != nil {
		return nil, fmt.Errorf("server.NewServer: %v", err)
	}

	return s, nil
}
