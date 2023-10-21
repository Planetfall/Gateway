package runserver

import (
	"fmt"

	"github.com/planetfall/framework/pkg/server"
	"github.com/planetfall/gateway/internal/service"
	"github.com/spf13/viper"
)

func getService(srv *server.Server) (*service.Service, error) {
	insecure := viper.GetBool(insecureFlag)
	projectID := viper.GetString(projectFlag)
	svc, err := service.NewService(service.ServiceOptions{
		Srv:       srv,
		ProjectID: projectID,
		Insecure:  insecure,
	})
	if err != nil {
		return nil, fmt.Errorf("service.NewService: %v", err)
	}

	return svc, nil
}
