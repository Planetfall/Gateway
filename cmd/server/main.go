package main

import (
	"fmt"
	"log"

	"github.com/Dadard29/planetfall/gateway/internal/server"
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
	env := flag.String("env", server.Production, "server environment")
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

func main() {
	setConfig()
	var connCfgList server.ConnectionConfigList
	err := viper.UnmarshalKey("services", &connCfgList)
	if err != nil {
		log.Fatalf("failed to load connections config: %v\n", err)
	}
	var connNameList []string
	for connName, _ := range connCfgList {
		connNameList = append(connNameList, connName)
	}
	log.Printf("loaded connections %v\n", connNameList)

	server, err := server.NewServer(
		viper.GetString("env"),
		viper.GetString("K_SERVICE"),
		viper.GetString("PORT"),
		connCfgList,
	)
	if err != nil {
		log.Fatalf("failed creating the server: %v\n", err)
	}

	server.Start()
}
