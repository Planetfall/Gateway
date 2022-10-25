package server

import (
	"context"
	"log"
	"net/http"
	"time"

	musicResearcherPb "github.com/Dadard29/planetfall/gateway/pkg/musicresearcher"
	"github.com/gin-gonic/gin"
)

const host = "music-researcher-twecq3u42q-ew.a.run.app:443"
const audience = "https://music-researcher-twecq3u42q-ew.a.run.app"

type searchParams struct {
	query string `form:"q"`
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

	ctx, err := s.conns.musicResearcher.getAuthenticatedCtx(ctx)
	if err != nil {
		log.Println("failed getting authenticated ctx")
		s.internalError(err, c)
		return
	}

	results, err := s.cls.musicResearcher.Search(ctx, &musicResearcherPb.Parameters{
		Query:        sp.query,
		GenreFilters: []string{},
		Limit:        10,
	})
	if err != nil {
		s.internalError(err, c)
		return
	}

	c.JSON(http.StatusOK, &results)
}
