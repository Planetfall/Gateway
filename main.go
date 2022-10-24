package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
	_ "github.com/Dadard29/planetfall/musicresearcher"
	"github.com/gin-gonic/gin"
)

type SearchParams struct {
	Query string `form:"q"`
}

var errorReporting *errorreporting.Client

func main() {
	ctx := context.Background()

	// read config from env
	serviceName := os.Getenv("K_SERVICE")

	log.Println("initializing metadata client...")
	metadataClient := metadata.NewClient(&http.Client{})
	projectID, err := metadataClient.ProjectID()
	if err != nil {
		log.Fatalf("metadata.ProjectID: %v\n", err)
	}

	log.Println("initializing error reporting")
	errorReporting, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: serviceName,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		log.Fatalf("errorreporting.NewClient: %v\n", err)
	}

	r := gin.Default()
	r.GET("/music-researcher/search", musicSearchController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to %s", port)
	}

	log.Printf("listening on port %s", port)

	r.Run(":" + port)
}
