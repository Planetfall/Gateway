package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/planetfall/gateway/docs"
	c "github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/download"
	"github.com/planetfall/gateway/internal/controller/search"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine *gin.Engine

	musicResearcherController *search.MusicResearcherController
	downloadController        *download.DownloadController
}

type RouterOptions struct {
	ProjectID           string
	Insecure            bool
	ErrorReportCallback func(err error)
	ConfigMap           ServicesConfig
}

const baseUrl = "/"

// supported services
const musicResearcherName = "music-researcher"
const youtubeDlJobName = "youtube-dl-job"

func NewRouter(opt RouterOptions) (*Router, error) {
	g := gin.Default()

	// cors
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	g.Use(cors.New(config))

	//music researcher
	mr, err := setupMusicResearcher(g, opt)
	if err != nil {
		return nil, fmt.Errorf("setupMusicResearcher: %v", err)
	}

	// youtube-dl job
	yt, err := setupDownload(g, opt)
	if err != nil {
		return nil, fmt.Errorf("setupDownload: %v", err)
	}

	g.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// g.Use(corsMiddleware())

	return &Router{
		engine: g,

		musicResearcherController: mr,
		downloadController:        yt,
	}, nil
}

func setupMusicResearcher(
	g *gin.Engine, opt RouterOptions) (*search.MusicResearcherController, error) {

	// music researcher
	conf, exists := opt.ConfigMap[musicResearcherName]
	if !exists {
		return nil, nil
	}

	log.Printf("setting up %s controller\n", musicResearcherName)

	mr, err := search.NewMusicResearcherController(
		c.ControllerOptions{
			Host:                conf.Host,
			Audience:            conf.Audience,
			Insecure:            opt.Insecure,
			ErrorReportCallback: opt.ErrorReportCallback,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"controller.NewMusicResearcherController: %v", err)
	}

	musicResearcherGroup := g.Group("/music-researcher")
	{
		musicResearcherGroup.GET("/search", mr.Search)
	}

	return mr, nil
}

func setupDownload(
	g *gin.Engine, opt RouterOptions) (*download.DownloadController, error) {

	conf, exists := opt.ConfigMap[youtubeDlJobName]
	if !exists {
		return nil, nil
	}
	log.Printf("setting up %s controller\n", youtubeDlJobName)
	yt, err := download.NewDownloadController(
		download.DownloadControllerOptions{

			ControllerOptions: c.ControllerOptions{
				Host:                conf.Host,
				ErrorReportCallback: opt.ErrorReportCallback,
			},

			ProjectID:        opt.ProjectID,
			LocationID:       conf.LocationID,
			QueueID:          conf.QueueID,
			SubscriptionID:   conf.SubscriptionID,
			WebsocketOrigins: conf.WebsocketOrigins,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("controllers.NewDownloadController: %v", err)
	}

	downloadGroup := g.Group("/download")
	{
		downloadGroup.GET("/url", yt.DownloadJob)
	}

	return yt, nil
}

// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
func (r *Router) Run(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: r.engine,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// blocking until quit is triggered
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Shutdown: %v", err)
	}

	log.Println("server closed")
	return nil
}

func (r *Router) Close() error {
	if err := r.musicResearcherController.Close(); err != nil {
		return fmt.Errorf("musicResearcher.Close: %v", err)
	}

	if err := r.downloadController.Close(); err != nil {
		return fmt.Errorf("download.Closed: %v", err)
	}

	return nil
}
