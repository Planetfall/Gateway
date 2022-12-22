package search

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/connection"
	"github.com/planetfall/gateway/internal/controller"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

type MusicResearcherController struct {
	controller.Controller

	conn   *connection.Connection
	client pb.MusicResearcherClient
}

func NewMusicResearcherController(
	opt controller.ControllerOptions) (*MusicResearcherController, error) {

	ctrl := controller.Controller{
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
// @Param       q     query string   true "Main user query"
// @Param       genre query []string true "Genre list"
// @Param       limit query int      true "Limit result count"
// @Success     200
// @Router      /music-researcher/search [get]
func (c *MusicResearcherController) Search(g *gin.Context) {

	// params
	var sp searchParam
	if err := g.ShouldBind(&sp); err != nil {
		c.BadRequest(fmt.Errorf("gin.ShouldBind: %v", err), g)
		return
	}

	ctx, cancel := c.GetContext(controller.DefaultTimeout)
	defer cancel()

	ctx, err := c.conn.AuthenticateContext(ctx)
	if err != nil {
		c.InternalError(fmt.Errorf("connection.AuthenticateContext: %v", err), g)
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
		c.InternalError(fmt.Errorf("musicResearcherClient.Search: %v", err), g)
		return
	}

	g.JSON(http.StatusOK, &results)
}
