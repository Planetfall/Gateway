package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	_ "github.com/Dadard29/planetfall/musicresearcher"
	pb "github.com/Dadard29/planetfall/musicresearcher"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
)

type SearchParams struct {
	Query string `form:"q"`
}

const host = "music-researcher-twecq3u42q-ew.a.run.app:443"
const audience = "https://music-researcher-twecq3u42q-ew.a.run.app"

func musicSearch(query string) {
	var opts []grpc.DialOption
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		log.Println("failed getting certs")
		log.Fatal(err)
	}

	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create auth context with token
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		log.Fatalf("idtoken.NewTokenSource: %v", err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("tokenSource.Token: %v", err)
	}

	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	c := pb.NewMusicResearcherClient(conn)

	r, err := c.Search(ctx, &pb.Parameters{
		Query:        query,
		GenreFilters: []string{},
		Limit:        10,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(r)
	log.Println(len(r.Tracks))
}

func main() {
	r := gin.Default()
	r.GET("/music-researcher/search", func(c *gin.Context) {
		var searchParams SearchParams
		if err := (c.ShouldBind(&searchParams)); err == nil {
			musicSearch(searchParams.Query)
		} else {
			log.Printf("failed to parse query: %v", err)
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
