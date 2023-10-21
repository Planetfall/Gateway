package runserver

import (
	"log"

	"github.com/spf13/viper"
)

func RunServer() {
	log.SetPrefix("[RUNSERVER] ")

	// config
	log.Printf("setting up the config")
	cfg, err := getConfig()
	if err != nil {
		log.Fatalf("getConfig: %v", err)
	}

	// server
	log.Printf("setting up the server")
	srv, err := getServer(cfg)
	if err != nil {
		log.Fatalf("getServer: %v", err)
	}

	// service
	log.Printf("setting up the service")
	svc, err := getService(srv)
	if err != nil {
		log.Fatalf("getService: %v", err)
	}

	port := viper.GetString(portFlag)
	if err := svc.Start(port); err != nil {
		log.Fatalf("svc.Start: %v", err)
	}

	log.Println("service stopped")
}
