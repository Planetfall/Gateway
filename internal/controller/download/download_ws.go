package download

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// @Description Payload to send to the download job
type urlDownloadParam struct {
	Url  string `json:"url"` // the youtube url to use for download
	Meta struct {
		Artist string `json:"artist"` // the artist
		Album  string `json:"album"`  // the album
		Track  string `json:"track"`  // the music track
	} `json:"meta"` // metadata infos to format the file
}

type urlDownloadPayload struct {
	JobKey string `json:"job_key"`

	urlDownloadParam
}

// @Summary     Download and save a new music file
// @Description Execute the Youtube-DL job using Cloud Task
// @Accept      json
// @Produces    json
// @Param       payload body urlDownloadParam true "Parameters to send to job"
// @Success     201
// @Router      /download/url [get]
func (c *DownloadController) DownloadJob(g *gin.Context) {

	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
	}

	ws, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		c.InternalError(fmt.Errorf("websocket.Upgrade: %v", err), g)
		return
	}

	log.Println("upgraded to websocket")

	if err := c.sub.RegisterWebsocket(ws); err != nil {
		c.InternalError(fmt.Errorf("subscriber.RegisterWebsocket: %v", err), g)
		return
	}

	// remove from stored sockets when connection ended
	defer c.sub.UnregisterWebsocket(ws)
	defer ws.Close()

	for {
		// main loop
		err := c.downloadJobWs(ws)

		// check if ws closed
		if websocket.IsCloseError(
			err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
			log.Println("closed websocket")
			break
		}

		if err != nil {
			log.Println(fmt.Errorf("download.downloadJobWs: %v", err))
		}
	}
}

func (c *DownloadController) downloadJobWs(ws *websocket.Conn) error {

	var dParam urlDownloadParam
	if err := ws.ReadJSON(&dParam); err != nil {
		return err
	}

	// add the new job to the current ws
	jKey, err := c.sub.AddNewJob(ws)
	if err != nil {
		return fmt.Errorf("subscriber.AddNewJob: %v", err)
	}

	// create task
	payload := urlDownloadPayload{
		urlDownloadParam: dParam,
		JobKey:           string(jKey),
	}
	createdTask, err := c.createTask(payload)
	if err != nil {
		return fmt.Errorf("download.createTask: %v", err)
	}

	// send back the created task
	if err := ws.WriteJSON(&createdTask); err != nil {
		return fmt.Errorf("websocket.WriteJSON: %v", err)
	}

	log.Println("created task, waiting for notif...")
	return nil
}
