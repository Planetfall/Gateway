package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/connection"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

type MusicResearcherController struct {
	Controller

	conn   *connection.Connection
	client pb.MusicResearcherClient
}

func NewMusicResearcherController(
	opt ControllerOptions) (*MusicResearcherController, error) {

	ctrl := Controller{
		opt.ErrorReportCallback,
	}

	conn, err := connection.NewConnection(opt.Insecure, opt.Host, opt.Audience)
	if err != nil {
		return nil, fmt.Errorf("connection.NewConnection: %v", err)
	}

	client := pb.NewMusicResearcherClient(conn.GrpcConn())

	return &MusicResearcherController{
		ctrl,
		conn,
		client,
	}, nil
}

func (c *MusicResearcherController) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("connection.Close: %v", err)
	}
	return nil
}

type searchParam struct {
	Query     string   `form:"q"`
	GenreList []string `form:"genre"`
	Limit     int      `form:"limit"`
}

// @Summary     Music search
// @Description Searchs for music in Spotify API
// @Accept      json
// @Produces    json
// @Param       q     query    string   true "Main user query"
// @Param       genre query    []string true "Genre list"
// @Param       limit query    int      true "Limit result count"
// @Success     200   {object} pb.Results
// @Router      /music-researcher/search [get]
func (c *MusicResearcherController) Search(g *gin.Context) {

	// params
	var sp searchParam
	if err := g.ShouldBind(&sp); err != nil {
		c.badRequest(fmt.Errorf("gin.ShouldBind: %v", err), g)
		return
	}

	ctx, cancel := c.getContext(defaultTimeout)
	defer cancel()

	ctx, err := c.conn.AuthenticateContext(ctx)
	if err != nil {
		c.internalError(fmt.Errorf("connection.AuthenticateContext: %v", err), g)
		return
	}

	results, err := c.client.Search(
		ctx,
		&pb.Parameters{
			Query:        sp.Query,
			GenreFilters: sp.GenreList,
			Limit:        int32(sp.Limit),
		},
	)
	if err != nil {
		c.internalError(fmt.Errorf("musicResearcherClient.Search: %v", err), g)
		return
	}

	g.JSON(http.StatusOK, &results)
}
