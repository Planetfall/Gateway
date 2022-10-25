package main

import (
	"log"

	"github.com/Dadard29/planetfall/gateway/internal/server"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("K_SERVICE", "gateway")

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed reading config file: %v\n", err)
	}

	server, err := server.NewServer(
		viper.GetString("K_SERVICE"),
		viper.GetString("PORT"),
	)
	if err != nil {
		log.Fatalf("failed creating the server: %v\n", err)
	}

	server.Start()
}
