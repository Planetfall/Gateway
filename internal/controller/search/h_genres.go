package search

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

// GetGenreList uses the SearchController client to interact with the
// dedicated service.
//
//	@Summary		List genres
//	@Description	List available genre in Spotify API
//	@Accept			json
//	@Produces		json
//	@Success		200
//	@Router			/music-researcher/genres [get]
func (c *SearchController) GetGenreList(g *gin.Context) {

	// get authentication context
	ctx, cancel := c.GetContext()
	defer cancel()

	ctx, err := c.conn.AuthenticateContext(ctx)
	if err != nil {
		c.InternalError(
			fmt.Errorf("connection.AuthenticateContext: %v", err), g)
		return
	}

	c.Logger.Printf("getting genre list")
	results, err := c.client.GetGenreList(ctx, &pb.Empty{})
	if err != nil {
		c.InternalError(fmt.Errorf("client.GetGenreList: %v", err), g)
		return
	}

	c.Logger.Printf("go %v genres", len(results.Genres))
	g.JSON(http.StatusOK, &results)
}
