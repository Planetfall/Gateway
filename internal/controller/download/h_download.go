package download

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/planetfall/gateway/internal/controller/download/task"
	"github.com/planetfall/gateway/internal/controller/download/websocket"
)

// Download upgrades HTTP request to a websocket
//
//	@Summary		Download and save a new music file
//	@Description	Execute the Youtube-DL job using Cloud Task
//	@Accept			json
//	@Produces		json
//	@Param			payload	body	task.Payload	true	"Parameters to send to job"
//	@Success		201
//	@Router			/download/url [get]
func (c *DownloadController) Download(g *gin.Context) {

	conn, err := c.websocket.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		c.BadRequest(fmt.Errorf("websocket.Upgrade: %v", err), g)
		return
	}

	c.Logger.Println("upgraded to websocket")

	// register websocket
	if err := c.websocketStore.Register(conn); err != nil {
		c.InternalError(fmt.Errorf("store.Register: %v", err), g)
		return
	}

	defer func() {

		// unregister websocket
		if err := c.websocketStore.Unregister(conn); err != nil {
			c.Logger.Printf("store.Unregister: %v", err)
		}

		// close websocket
		if err := conn.Close(); err != nil {
			c.Logger.Printf("websocket.Close: %v", err)
		}
	}()

	for {
		// main loop
		err := c.Loop(conn)

		// check if ws closed
		if c.websocket.IsClosed(err) {
			c.Logger.Println("closed websocket")
			break
		}

		if err != nil {
			c.Logger.Printf("download.Loop: %v", err)
			continue
		}
	}
}

func (c *DownloadController) Loop(conn websocket.Conn) error {

	// read current JSON
	var payload task.Payload
	if err := conn.ReadJSON(&payload); err != nil {
		return fmt.Errorf("ws.ReadJSON: %v", err)
	}

	jKey, err := c.websocketStore.AddNewJob(conn)
	if err != nil {
		return fmt.Errorf("subscriber.AddNewJob: %v", err)
	}

	// create task
	taskPayload := task.TaskPayload{
		Payload: payload,
		JobKey:  string(jKey),
	}
	createdTask, err := c.taskClient.CreateTask(taskPayload)
	if err != nil {
		return fmt.Errorf("download.createTask: %v", err)
	}

	// send back the created task
	if err := conn.WriteJSON(&payload); err != nil {
		return fmt.Errorf("websocket.WriteJSON: %v", err)
	}

	c.Logger.Printf("created task %s", createdTask.Name)
	return nil
}

// ReceiveCallback is called when a message is received from the subscription.
// It ensures that the message is acknowledged.
// The received message is parsed. Then, the calling websocket is retrieved in
// the store using the job ordering key. The parsed message is written on this
// websocket.
func (c *DownloadController) OnReceive(
	ctx context.Context, message *pubsub.Message) {

	defer message.Ack()

	// parse pMsg content
	jobStatus, err := c.sub.NewJobStatus(message)
	if err != nil {
		c.Logger.Println(fmt.Errorf("subscriber.NewJobStatus: %v", err))
		return
	}

	c.Logger.Printf("Received job status with key: %s | code: %d",
		jobStatus.OrderingKey, jobStatus.Code)

	// retrieve the websocket using the message ordering key
	orderingKey := websocket.Key(jobStatus.OrderingKey)
	ws, err := c.websocketStore.GetWebsocket(orderingKey)
	if err != nil {
		c.Logger.Println(fmt.Errorf("store.GetWebsocket: %v", err))
		return
	}

	// notify to ws
	if err := ws.WriteJSON(&jobStatus); err != nil {
		c.Logger.Println(fmt.Errorf("websocket.WriteJSON: %v", err))
		return
	}
}
