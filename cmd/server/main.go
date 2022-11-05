package main

import (
	"fmt"
	"log"

	"github.com/planetfall/gateway/pkg/config"
	rs "github.com/planetfall/gateway/pkg/runserver"
	"github.com/spf13/viper"
)

func getServiceConfigMap() (config.ServiceConfigMap, error) {
	var serviceConfigMap config.ServiceConfigMap
	err := viper.UnmarshalKey("services", &serviceConfigMap)
	if err != nil {
		return config.ServiceConfigMap{}, fmt.Errorf("viper.UnmarshalKey: %v", err)
	}

	return serviceConfigMap, nil
}

func main() {
	config.ReadConfig()
	serviceConfigMap, err := getServiceConfigMap()
	if err != nil {
		log.Fatalf("failed to get connection list from config: %v", err)
	}

	err = rs.RunServer(
		rs.RunServerOptions{
			ServiceName: viper.GetString("K_SERVICE"),
			ProjectID:   viper.GetString("project-id"),
			Port:        viper.GetString("PORT"),
			Insecure:    viper.GetBool("insecure"),
			SvcConfig:   serviceConfigMap,

			Environment: viper.GetString("env"),
		},
	)
	if err != nil {
		log.Fatalf("runserver.RunServer: %v\n", err)
	}
}
