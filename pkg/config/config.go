package config

import (
	"fmt"
	"log"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

const (
	defaultPort    = "8080"
	defaultService = "gateway"

	defaultEnv        = Production
	defaultInsecure   = false
	defaultConfigFile = "./config/config.yaml"
)

// Read config from environment, commandline and file
// Reference:
// - PORT				(environment):	set the HTTP port to listen too (default: 8080)
// - K_SERVICE	(environment):	set with the current Cloud Run service name,
// included in Cloud Run environment:
// https://cloud.google.com/run/docs/container-contract#env-vars

// - env				(commandline):	set with environments values
// - insecure		(commandline):	boolean, use to interact with others services
// without authentication, DO NOT USE IN PRODUCTION
// - config			(commandline):	set the config file path

// - configfile									set project ID and services connections
func ReadConfig() {
	// from env
	viper.SetDefault("PORT", defaultPort)
	viper.SetDefault("K_SERVICE", defaultService)

	viper.BindEnv("PORT")
	viper.BindEnv("K_SERVICE")

	// from cmd line
	flag.String("env", defaultEnv, "server environment")
	flag.Bool("insecure", defaultInsecure, "insecure connection with microservices if true")
	flag.String("config", defaultConfigFile, "config file path")
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// from config file
	configFile := viper.GetString("config")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed reading config file: %v\n", err)
	}
	log.Printf("loaded config file %s...\n", configFile)
}

func getServiceConfigMap() (*ServiceConfigMap, error) {
	var serviceConfigMap ServiceConfigMap
	err := viper.UnmarshalKey("services", &serviceConfigMap)
	if err != nil {
		return nil, fmt.Errorf("viper.UnmarshalKey: %v", err)
	}

	return &serviceConfigMap, nil
}
