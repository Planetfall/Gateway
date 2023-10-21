package service

import (
	"log"
	"reflect"
)

// Keys used to retrieve controller configuration and set the routes
const (
	MusicResearcher = "music-researcher"
	Downloader      = "downloader"
)

// controllerConfig holds the basic configuration needed for a GRPC controller
type controllerConfig struct {
	Target string `mapstructure:"target" validate:"required"`
}

// downloadControllerConfig holds specific configuration for the download
// controller
type downloadControllerConfig struct {
	Target string `mapstructure:"target" validate:"required"`

	LocationID     string   `mapstructure:"location" validate:"required"`
	QueueID        string   `mapstructure:"queue" validate:"required"`
	SubscriptionID string   `mapstructure:"subscription" validate:"required"`
	Origins        []string `mapstructure:"origins" validate:"required"`
}

func logConfig(logger *log.Logger, cfg interface{}) {
	v := reflect.ValueOf(cfg)
	t := v.Type()

	logger.Println("--------------")
	for i := 0; i < v.NumField(); i++ {
		logger.Printf("%s:\t %v", t.Field(i).Name, v.Field(i).Interface())
	}
	logger.Println("--------------")
}
