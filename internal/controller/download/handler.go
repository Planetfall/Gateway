package download

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// download payload contains the needed fields to perform the download.
type downloadPayload struct {
	Url  string `json:"url"` // the youtube url to use for download
	Meta struct {
		Artist string `json:"artist"` // the artist
		Album  string `json:"album"`  // the album
		Track  string `json:"track"`  // the music track
	} `json:"meta"` // metadata infos to format the file
}

// taskDownloadPaylaod is the payload sent to the download job.
type taskDownloadPayload struct {
	JobKey string `json:"job_key"`
	downloadPayload
}

// checkOrigins checks if the Origin header of the HTTP request is allowed
// by looking into the controller setup origins.
func (c *DownloadController) checkOrigins(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	// fixme
	// no origin to verify (aka, not called by javascript)
	if origin == "" {
		return true
	}

	for _, authorizedOrigin := range c.origins {
		if authorizedOrigin == origin {
			return true
		}
	}
	return false
}

// Download upgrades HTTP request to a websocket
//
//	@Summary		Download and save a new music file
//	@Description	Execute the Youtube-DL job using Cloud Task
//	@Accept			json
//	@Produces		json
//	@Param			payload	body	downloadPayload	true	"Parameters to send to job"
//	@Success		201
//	@Router			/download/url [get]
func (c *DownloadController) Download(g *gin.Context) {

	ws, err := c.upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		c.InternalError(fmt.Errorf("websocket.Upgrade: %v", err), g)
		return
	}

	c.Logger.Println("upgraded to websocket")

	if err := c.sub.RegisterWebsocket(ws); err != nil {
		c.InternalError(fmt.Errorf("subscriber.RegisterWebsocket: %v", err), g)
		return
	}

	// remove from stored sockets when connection ended
	defer c.sub.UnregisterWebsocket(ws)
	defer ws.Close()

	for {
		// main loop
		err := c.downloadFromWebsocket(ws)

		// check if ws closed
		if websocket.IsCloseError(
			err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
			c.Logger.Println("closed websocket")
			break
		}

		if err != nil {
			c.Logger.Println(fmt.Errorf("download.downloadJobWs: %v", err))
		}
	}
}

// downloadFromWebsocket takes a websocket and setup a new download job.
func (c *DownloadController) downloadFromWebsocket(ws *websocket.Conn) error {

	// read current JSON
	var dParam downloadPayload
	if err := ws.ReadJSON(&dParam); err != nil {
		return fmt.Errorf("ws.ReadJSON: %v", err)
	}

	// add the new job to the current ws
	jKey, err := c.sub.AddNewJob(ws)
	if err != nil {
		return fmt.Errorf("subscriber.AddNewJob: %v", err)
	}

	// create task
	payload := taskDownloadPayload{
		downloadPayload: dParam,
		JobKey:          string(jKey),
	}
	createdTask, err := c.createTask(payload)
	if err != nil {
		return fmt.Errorf("download.createTask: %v", err)
	}

	// send back the created task
	if err := ws.WriteJSON(&payload); err != nil {
		return fmt.Errorf("websocket.WriteJSON: %v", err)
	}

	c.Logger.Printf("created task %v, waiting for notif...", createdTask.Name)
	return nil
}
