package router

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	c "github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/pkg/config"
)

const baseUrl = "/"

// supported services
const musicResearcherName = "music-researcher"

type Router struct {
	engine *gin.Engine
}

func formatRoute(service string, endpoint string) string {
	return fmt.Sprintf("%s%s/%s", baseUrl, service, endpoint)
}

func NewRouter(
	configMap config.ServiceConfigMap,
	errorReportCallback func(err error),
	insecure bool,
) (*Router, error) {

	g := gin.Default()
	if conf, exists := configMap[musicResearcherName]; exists {
		log.Printf("setting up %s controller\n", musicResearcherName)

		mr, err := c.NewMusicResearcherController(
			c.ControllerOptions{
				Host:                conf.Host,
				Audience:            conf.Audience,
				Insecure:            insecure,
				ErrorReportCallback: errorReportCallback,
			})
		if err != nil {
			return nil, fmt.Errorf("controller.NewMusicResearcherController: %v", err)
		}

		g.GET(formatRoute(musicResearcherName, "search"), mr.Search)
	}

	return &Router{
		engine: g,
	}, nil
}

func (r *Router) Run(addr ...string) error {
	return r.engine.Run(addr...)
}
