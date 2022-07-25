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
	r.GET("/music-researcher/search", func(c *gin.Context) {
		var searchParams SearchParams
		if err := (c.ShouldBind(&searchParams)); err == nil {
			log.Println(searchParams.Query)
		} else {
			log.Println(err)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to %s", port)
	}

	log.Printf("listening on port %s", port)

	r.Run(":" + port)
}
