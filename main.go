package main

import (
	"log"
	"os"

	_ "github.com/Dadard29/planetfall/musicresearcher"
	"github.com/gin-gonic/gin"
)

type SearchParams struct {
	Query string `form:"q"`
}

func main() {
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
