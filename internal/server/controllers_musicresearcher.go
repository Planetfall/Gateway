package server

import (
	"context"
	"log"
	"net/http"
	"time"

	musicResearcherPb "github.com/Dadard29/planetfall/gateway/pkg/musicresearcher"
	"github.com/gin-gonic/gin"
)

type searchParams struct {
	Query string `form:"q"`
}

func (s *Server) musicSearch(c *gin.Context) {
	var sp searchParams
	if err := (c.ShouldBind(&sp)); err != nil {
		log.Println("failed parsing search params")
		s.badRequest(err, c)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, err := s.conns[conn_MusicResearcher].getAuthenticatedCtx(ctx)
	if err != nil {
		log.Println("musicResearcher.getAuthenticatedCtx: %v\n", err)
		s.internalError(err, c)
		return
	}

	results, err := s.cls.musicResearcher.Search(ctx, &musicResearcherPb.Parameters{
		Query:        sp.Query,
		GenreFilters: []string{},
		Limit:        10,
	})
	if err != nil {
		s.internalError(err, c)
		return
	}

	c.JSON(http.StatusOK, &results)
}
