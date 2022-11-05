package router

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/planetfall/gateway/docs"
	c "github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/pkg/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine *gin.Engine
}

type RouterOptions struct {
	ProjectID           string
	Insecure            bool
	ErrorReportCallback func(err error)
	ConfigMap           config.ServiceConfigMap
}

const baseUrl = "/"

// supported services
const musicResearcherName = "music-researcher"
const youtubeDlJobName = "youtube-dl-job"

func formatRoute(service string, endpoint string) string {
	return fmt.Sprintf("%s%s/%s", baseUrl, service, endpoint)
}

func NewRouter(opt RouterOptions) (*Router, error) {

	g := gin.Default()

	// music researcher
	if conf, exists := opt.ConfigMap[musicResearcherName]; exists {
		log.Printf("setting up %s controller\n", musicResearcherName)

		mr, err := c.NewMusicResearcherController(
			c.ControllerOptions{
				Host:                conf.Host,
				Audience:            conf.Audience,
				Insecure:            opt.Insecure,
				ErrorReportCallback: opt.ErrorReportCallback,
			})
		if err != nil {
			return nil, fmt.Errorf("controller.NewMusicResearcherController: %v", err)
		}

		musicResearcherGroup := g.Group("/music-researcher")
		{
			musicResearcherGroup.GET("/search", mr.Search)
		}
	}

	// youtube-dl job
	if conf, exists := opt.ConfigMap[youtubeDlJobName]; exists {
		log.Printf("setting up %s controller\n", youtubeDlJobName)
		yt, err := c.NewDownloadController(
			c.DownloadControllerOptions{
				ProjectID:  opt.ProjectID,
				LocationID: conf.LocationID,
				QueueID:    conf.QueueID,
			},
			c.ControllerOptions{
				Host:                conf.Host,
				ErrorReportCallback: opt.ErrorReportCallback,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("controllers.NewDownloadController: %v", err)
		}

		downloadGroup := g.Group("/download")
		{
			downloadGroup.POST("/url", yt.DownloadJob)
		}
	}

	g.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Router{
		engine: g,
	}, nil
}

func (r *Router) Run(addr ...string) error {
	return r.engine.Run(addr...)
}
