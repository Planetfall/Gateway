package search

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

// searchParameters holds the arguments passed to the Search endpoint
type searchParameters struct {
	Query     string   `form:"q"`
	GenreList []string `form:"genre"`
	Limit     int      `form:"limit"`
}

// Search uses the SearchController client to interact with the
// dedicated service.
//
//	@Summary		Music search
//	@Description	Searchs for music in Spotify API
//	@Accept			json
//	@Produces		json
//	@Param			q		query	string		true	"Main user query"
//	@Param			genre	query	[]string	true	"Genre list"
//	@Param			limit	query	int			true	"Limit result count"
//	@Success		200
//	@Router			/music-researcher/search [get]
func (c *SearchController) Search(g *gin.Context) {

	// parse the parameters
	var sp searchParameters
	if err := g.ShouldBind(&sp); err != nil {
		c.BadRequest(fmt.Errorf("gin.ShouldBind: %v", err), g)
		return
	}

	// get authentication context
	ctx, cancel := c.GetContext()
	defer cancel()

	ctx, err := c.conn.AuthenticateContext(ctx)
	if err != nil {
		c.InternalError(
			fmt.Errorf("connection.AuthenticateContext: %v", err), g)
		return
	}

	// use the client to perform the search
	c.Logger.Printf("searching with query: `%v` | genres: `%v` | limit: %v",
		sp.Query, sp.GenreList, sp.Limit)
	results, err := c.client.Search(
		ctx,
		&pb.Parameters{
			Query:        sp.Query,
			GenreFilters: sp.GenreList,
			Limit:        int32(sp.Limit),
		},
	)
	if err != nil {
		c.InternalError(
			fmt.Errorf("client.Search: %v", err), g)
		return
	}

	// send back the result
	c.Logger.Printf("searched and got %v tracks", len(results.Tracks))
	g.JSON(http.StatusOK, &results)
}
