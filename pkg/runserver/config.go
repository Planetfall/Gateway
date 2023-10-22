package runserver

import (
	"fmt"
	"log"

	"github.com/planetfall/framework/pkg/config"
	"github.com/spf13/viper"
)

const (
	portFlag     = "port"
	serviceFlag  = "service"
	projectFlag  = "project"
	insecureFlag = "insecure"
)

func getConfig() (config.Config, error) {
	entries := []config.Entry{
		{
			Flag:         portFlag,
			DefaultValue: "8080",
			Description:  "the exposed port of the service",
			EnvKey:       "PORT",
		},
		{
			Flag:         serviceFlag,
			DefaultValue: "cloud-microservice",
			Description:  "the service name",
			EnvKey:       "K_SERVICE",
		},
		{
			Flag:         insecureFlag,
			DefaultValue: "false",
			Description:  "sets the connection to services to be insecure",
			EnvKey:       "INSECURE",
		},
		{
			Flag:         projectFlag,
			DefaultValue: "project-id",
			Description:  "the service name",
			EnvKey:       "K_SERVICE",
		},
	}

	cfg, err := config.NewConfig(entries)
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: %v", err)
	}

	log.Println("--------------")
	for _, entry := range entries {
		key := entry.Flag
		value := viper.GetString(key)
		log.Printf("- %s \t %s", key, value)
	}
	log.Println("--------------")

	return cfg, nil
}
