package service

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/download"
	"github.com/planetfall/gateway/internal/controller/search"
	"github.com/spf13/viper"
)

// controllerBuilders contains the controllers builder, associated with a key.
// The key is used to retrieve configuration from [viper] and builds the
// controller route.
var controllerBuilders = map[string]func(
	svcControllerOptions) (svcController, error){

	MusicResearcher: newSearchController,
	Downloader:      newDownloadController,
}

// svcController is the controller type used by the service.
type svcController interface {
	// Closes the controller and its resources.
	Close() error

	// Returns the the controller name.
	Name() string
}

// svcControllerOptions are the parameters for the controller builders.
type svcControllerOptions struct {
	// The key used to retrieve the controller configuration from [viper].
	cfgKey string

	// The logger for the controller
	logger *log.Logger

	// The router group with the base controller route set.
	group *gin.RouterGroup

	// The callback used by the controller when error is encountered in the
	// handlers.
	reportErrorCallback func(err error)

	// Used by the GRPC controllers
	insecure bool

	// Used to interact with Cloud features
	projectID string
}

// newSearchController creates a new SearchController
func newSearchController(opt svcControllerOptions) (svcController, error) {

	var cfg controllerConfig
	err := viper.UnmarshalKey(opt.cfgKey, &cfg)
	if err != nil {
		return nil, fmt.Errorf("viper.UnmarshalKey: %v", err)
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	logConfig(opt.logger, cfg)

	ctrlOpt := search.SearchControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:        opt.cfgKey,
			Target:      cfg.Target,
			ReportError: opt.reportErrorCallback,
			Logger:      opt.logger,
		},
		Insecure: opt.insecure,
	}
	ctrl, err := search.NewSearchController(ctrlOpt)
	if err != nil {
		return nil, fmt.Errorf("search.NewSearchController: %v", err)
	}

	opt.group.GET("/search", ctrl.Search)
	opt.group.GET("/genres", ctrl.GetGenreList)

	var svcCtrl svcController = ctrl
	return svcCtrl, nil
}

// newDownloadController creates a new DownloadController
func newDownloadController(opt svcControllerOptions) (svcController, error) {

	var cfg downloadControllerConfig
	err := viper.UnmarshalKey(opt.cfgKey, &cfg)
	if err != nil {
		return nil, fmt.Errorf("viper.UnmarshalKey: %v", err)
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	logConfig(opt.logger, cfg)

	ctrlOpt := download.DownloadControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:        opt.cfgKey,
			Target:      cfg.Target,
			ReportError: opt.reportErrorCallback,
			Logger:      opt.logger,
		},
		ProjectID:      opt.projectID,
		LocationID:     cfg.LocationID,
		QueueID:        cfg.QueueID,
		SubscriptionID: cfg.SubscriptionID,
		Origins:        cfg.Origins,
	}
	ctrl, err := download.NewDownloadController(ctrlOpt)
	if err != nil {
		return nil, fmt.Errorf("download.NewDownloadController: %v", err)
	}

	opt.group.GET("/url", ctrl.Download)

	var svcCtrl svcController = ctrl
	return svcCtrl, nil
}
