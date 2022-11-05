package main

import (
	"fmt"
	"log"

	"github.com/planetfall/gateway/internal/server"
	"github.com/planetfall/gateway/pkg/config"
	env "github.com/planetfall/gateway/pkg/environments"
	"github.com/planetfall/gateway/pkg/options"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func setConfig() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("K_SERVICE", "gateway")

	// from env
	viper.BindEnv("PORT")
	viper.BindEnv("K_SERVICE")

	// from cmd line
	env := flag.String("env", env.Production, "server environment")
	flag.Bool("insecure", false, "insecure connection with microservices if true")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// from config file
	configFile := fmt.Sprintf("./config/config.%s.yaml", *env)
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed reading config file: %v\n", err)
	}
	log.Printf("loaded config file %s...\n", configFile)
}

func getServiceConfigMap() (*config.ServiceConfigMap, error) {
	var serviceConfigMap config.ServiceConfigMap
	err := viper.UnmarshalKey("services", &serviceConfigMap)
	if err != nil {
		return nil, fmt.Errorf("viper.UnmarshalKey: %v", err)
		log.Fatalf("failed to load connections config: %v\n", err)
	}

	return &serviceConfigMap, nil
}

func main() {
	setConfig()
	serviceConfigMap, err := getServiceConfigMap()
	if err != nil {
		log.Fatalf("failed to get connection list from config: %v", err)
	}

	server, err := server.NewServer(
		options.ServerOptions{
			Env:          viper.GetString("env"),
			ServiceName:  viper.GetString("K_SERVICE"),
			ProjectID:    viper.GetString("project-id"),
			Port:         viper.GetString("PORT"),
			Insecure:     viper.GetBool("insecure"),
			SvcConfigMap: *serviceConfigMap,
		},
	)
	if err != nil {
		log.Fatalf("failed creating the server: %v\n", err)
	}

	server.Start()
}
