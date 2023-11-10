package download_test

import (
	"log"
	"testing"

	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/download"
	"github.com/stretchr/testify/assert"
)

func TestNewDownloadController(t *testing.T) {
	// given
	optGiven := download.DownloadControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:   "name",
			Target: "target",
			Logger: log.Default(),
		},
		ProjectID: "project-id",
	}

	// when
	c, err := download.NewDownloadController(optGiven)

	// then
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	// given
	optGiven := download.DownloadControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:   "name",
			Target: "target",
			Logger: log.Default(),
		},
		ProjectID: "project-id",
	}
	c, err := download.NewDownloadController(optGiven)
	assert.Nil(t, err)

	// when
	err = c.Close()
	assert.Nil(t, err)
}
