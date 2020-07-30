package impl

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cloudor-io/cloudctl/pkg/request"
	"github.com/gorilla/websocket"
)

const DefaultServerURL string = "https://cloudor.dev"

// JobStatus is used for streaming job status to the client
type JobStatus struct {
	UserName    string `json:"user_name,omitempty"`
	ID          string `json:"id,omitempty"`
	Status      string `json:"status,omitempty"`
	StatusCode  int    `json:"status_code,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Description string `json:"description,omitempty"`
}

func CheckingJob(jobMsg *request.RunJobMessage, username *string, token *string) (*request.RunJobMessage, error) {
	// serverURL := "http://localhost:3001" // DefaultServerURL
	serverURL := DefaultServerURL
	scheme := "ws"
	if strings.HasPrefix(serverURL, "http://") {
		serverURL = strings.Replace(serverURL, "http://", "", -1)
		scheme = "ws"
	}
	if strings.HasPrefix(serverURL, "https://") {
		serverURL = strings.Replace(serverURL, "https://", "", -1)
		scheme = "wss"
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: scheme, Host: serverURL, Path: "/ws/v1"}

	header := http.Header{}
	header.Add("From", *username)
	header.Add("Sec-WebSocket-Protocol", "Bearer "+*token)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("Error dial websocket %v:", err)
		return nil, err
	}
	defer c.Close()

	done := make(chan struct{})
	defer close(done)
	go func() {

		for {
			mtype, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if mtype != websocket.TextMessage {
				log.Printf("Received wrong msg type %d", mtype)
				return
			}
			job := JobStatus{}
			err = json.Unmarshal(message, &job)
			if err != nil {
				log.Printf("error parsing jobmsg: %v,", err)
				return
			}
			log.Printf("job status \"%s\": %s", job.Status, job.Description)
			if job.Status == "finished" || job.Status == "failed" || job.Status == "canceled" {
				done <- struct{}{}
				return
			}
		}
	}()

	jobBytes, _ := json.Marshal(jobMsg)
	err = c.WriteMessage(websocket.TextMessage, jobBytes)
	if err != nil {
		log.Printf("Error writing job to websocket: %v", err)
		return nil, err
	}
	doneFlag := false
	for {
		select {
		case <-done:
			doneFlag = true
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return nil, err
			}
			break
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close error:", err)
				return nil, err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil, errors.New("interrupted, quitting")
		}
		if doneFlag {
			break
		}
	}
	// get job msg
	resp, err := request.GetCloudor(username, token, "/job/user/"+*username+"/id/"+jobMsg.ID)
	if err != nil {
		log.Printf("getting job failed %v", err)
		return nil, err
	}
	// somewhere the json is encoded twice, unquote it TODO

	jobMessage := request.RunJobMessage{}
	err = json.Unmarshal(resp, &jobMessage)
	if err != nil {
		log.Printf("Internal error, cann't parse job response: %v", err)
		return nil, err
	}
	// log.Printf("Return job message %+v", jobMessage)
	return &jobMessage, nil
}
