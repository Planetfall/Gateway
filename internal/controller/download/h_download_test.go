package download_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/controller"
	"github.com/planetfall/gateway/internal/controller/download"
	"github.com/planetfall/gateway/internal/controller/download/mocks"
	"github.com/planetfall/gateway/internal/controller/download/subscriber"
	"github.com/planetfall/gateway/internal/controller/download/websocket"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {

	// given
	taskClientGiven := mocks.NewTaskClientMock().(*mocks.TaskClientMock)
	subscriberGiven := mocks.NewSubscriberMock().(*mocks.SubscriberMock)
	websocketGiven := mocks.NewWebsocketMock().(*mocks.WebsocketMock)
	storeGiven := websocket.NewStore()

	providerGiven := mocks.NewProviderMock(
		taskClientGiven,
		subscriberGiven,
		websocketGiven,
		storeGiven).(*mocks.ProviderMock)

	subscriberGiven.On("Listen").Return(nil)
	subscriberGiven.On("Close").Return(nil)

	// when
	addr := mocks.NewAddrMock("192.168.0.1")

	connGiven := mocks.NewConnMock().(*mocks.ConnMock)
	connGiven.On("ReadJSON").Return(nil)
	connGiven.On("WriteJSON").Return(nil)
	connGiven.On("RemoteAddr").Return(addr)
	connGiven.On("Close").Return(nil)

	websocketGiven.On("Upgrade").Return(connGiven, nil)
	websocketGiven.On("IsClosed").Return(true)

	taskClientGiven.On("CreateTask").Return(&cloudtaskspb.Task{Name: "newTask"}, nil)
	taskClientGiven.On("Close").Return(nil)

	// then
	opt := download.DownloadControllerOptions{
		Provider: providerGiven,
		ControllerOptions: controller.ControllerOptions{
			Logger: log.Default(),
		},
	}
	c, err := download.NewDownloadController(opt)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	// wait for listen
	time.Sleep(1 * time.Second)

	wGiven := httptest.NewRecorder()
	gGiven, _ := gin.CreateTestContext(wGiven)
	c.Download(gGiven)

	c.Close()

	taskClientGiven.AssertExpectations(t)
	subscriberGiven.AssertExpectations(t)
	websocketGiven.AssertExpectations(t)
	providerGiven.AssertExpectations(t)
}

func TestDownload_OnReceive(t *testing.T) {
	// given
	taskClientGiven := mocks.NewTaskClientMock().(*mocks.TaskClientMock)
	subscriberGiven := mocks.NewSubscriberMock().(*mocks.SubscriberMock)
	websocketGiven := mocks.NewWebsocketMock().(*mocks.WebsocketMock)
	storeGiven := websocket.NewStore()

	providerGiven := mocks.NewProviderMock(
		taskClientGiven,
		subscriberGiven,
		websocketGiven,
		storeGiven).(*mocks.ProviderMock)

	connGiven := mocks.NewConnMock().(*mocks.ConnMock)
	addr := mocks.NewAddrMock("192.168.0.1")
	connGiven.On("RemoteAddr").Return(addr)
	connGiven.On("WriteJSON").Return(nil)
	err := storeGiven.Register(connGiven)
	assert.Nil(t, err)

	key, err := storeGiven.AddNewJob(connGiven)
	assert.Nil(t, err)

	messageGiven := &pubsub.Message{}
	jobStatusGiven := &subscriber.JobStatus{
		OrderingKey: string(key),
		Code:        200,
	}
	subscriberGiven.On("NewJobStatus", messageGiven).Return(jobStatusGiven, nil)
	subscriberGiven.On("Listen").Return(nil)

	// then
	opt := download.DownloadControllerOptions{
		Provider: providerGiven,
		ControllerOptions: controller.ControllerOptions{
			Logger: log.Default(),
		},
	}
	c, err := download.NewDownloadController(opt)
	assert.Nil(t, err)
	assert.NotNil(t, c)

	// wait for listen
	time.Sleep(1 * time.Second)

	c.OnReceive(context.Background(), messageGiven)

	taskClientGiven.AssertExpectations(t)
	subscriberGiven.AssertExpectations(t)
	websocketGiven.AssertExpectations(t)
	providerGiven.AssertExpectations(t)
}

func TestDownload_withUpgradeError(t *testing.T) {

	wGiven := httptest.NewRecorder()
	gGiven, _ := gin.CreateTestContext(wGiven)
	req, err := http.NewRequest(http.MethodGet, "", nil)
	assert.Nil(t, err)
	gGiven.Request = req

	reportErrorGiven := func(err error) {
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Upgrade")
	}
	optGiven := download.DownloadControllerOptions{
		ControllerOptions: controller.ControllerOptions{
			Name:        "name",
			Target:      "target",
			Logger:      log.Default(),
			ReportError: reportErrorGiven,
		},
		ProjectID: "project-id",
	}

	c, err := download.NewDownloadController(optGiven)
	assert.Nil(t, err)

	c.Download(gGiven)
}
