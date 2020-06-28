package impl

import (
	"encoding/json"
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

func CheckingJob(jobMsg *request.RunJobMessage, username *string, token *string) {
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
	log.Printf("connecting to %s", u.String())

	header := http.Header{}
	header.Add("From", *username)
	header.Add("Sec-WebSocket-Protocol", "Bearer "+*token)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("Error dial websocket %v:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mtype, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv typs %d msg: %s", mtype, message)
		}
	}()

	// ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()
	jobBytes, _ := json.Marshal(jobMsg)

	err = c.WriteMessage(websocket.TextMessage, jobBytes)
	if err != nil {
		log.Printf("Error writing job to websocket: %v", err)
		return
	}
	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		return
		// 	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
