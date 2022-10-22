package main

import (
	"net/http"

	_ "github.com/Dadard29/planetfall/musicresearcher"
	pb "github.com/Dadard29/planetfall/musicresearcher"
	"github.com/gin-gonic/gin"
)

const host = "music-researcher-twecq3u42q-ew.a.run.app:443"
const audience = "https://music-researcher-twecq3u42q-ew.a.run.app"

func musicSearchController(c *gin.Context) {
	var searchParams SearchParams
	if err := (c.ShouldBind(&searchParams)); err != nil {
		badRequest(err, c)
		return
	}

	results, err := musicSearch(searchParams.Query)
	if err != nil {
		internalError(err, c)
		return
	}

	c.JSON(http.StatusOK, &results)
}

func musicSearch(query string) (*pb.Results, error) {

	conn, err := getGrpcConn(host, audience)
	if err != nil {
		return nil, err
	}

	ctx, err := getAuthenticatedCtx(conn)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	c := pb.NewMusicResearcherClient(conn)

	return c.Search(ctx, &pb.Parameters{
		Query:        query,
		GenreFilters: []string{},
		Limit:        10,
	})
}
