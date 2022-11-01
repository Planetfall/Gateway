package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/connection"
	pb "github.com/planetfall/gateway/pkg/musicresearcher"
)

type MusicResearcherController struct {
	Controller

	conn   *connection.Connection
	client pb.MusicResearcherClient
}

func NewMusicResearcherController(
	opt ControllerOptions) (*MusicResearcherController, error) {

	conn, err := connection.NewConnection(opt.Insecure, opt.Host, opt.Audience)
	if err != nil {
		return nil, fmt.Errorf("connection.NewConnection: %v", err)
	}

	client := pb.NewMusicResearcherClient(conn.GrpcConn())

	return &MusicResearcherController{
		Controller{opt.ErrorReportCallback},
		conn,
		client,
	}, nil
}

type searchParam struct {
	Query string `form:"q"`
}

// @Summary     Music search
// @Description Searchs for music in Spotify API
// @Accept      json
// @Produces    json
// @Param       q   query    string true "Main user query"
// @Success     200 {object} pb.Results
// @Router      /music-researcher/search [get]
func (c *MusicResearcherController) Search(g *gin.Context) {
	var sp searchParam
	if err := g.ShouldBind(&sp); err != nil {
		c.badRequest(fmt.Errorf("gin.Context.ShouldBind: %v", err), g)
		return
	}

	ctx, cancel := c.getContext(defaultTimeout)
	defer cancel()

	ctx, err := c.conn.AuthenticateContext(ctx)
	if err != nil {
		c.internalError(fmt.Errorf("connection.AuthenticateContext: %v", err), g)
		return
	}

	// fixme
	results, err := c.client.Search(
		ctx,
		&pb.Parameters{
			Query:        sp.Query,
			GenreFilters: []string{},
			Limit:        10,
		},
	)
	if err != nil {
		c.internalError(fmt.Errorf("musicResearcherClient.Search: %v", err), g)
		return
	}

	g.JSON(http.StatusOK, &results)
}
