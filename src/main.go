package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"adriftdev.com/eventual-agent/src/logging"
	"adriftdev.com/eventual-agent/src/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {

	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		pool.Logging.Warn((err.Error()))
	}

	client := &websocket.Client{
		ID:   r.RemoteAddr,
		Conn: conn,
		Pool: pool,
	}

	client.Read()
}

func handlePublishPostRequest(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		pool.Logging.Error(err.Error())
	}
	//Convert the body to type string
	sb := string(body)
	pool.Logging.Debug(sb)
	pool.Publish <- websocket.PublishRequest{}
	w.WriteHeader(202)
}

func setupRoutes(l *logging.Logger) {
	pool := websocket.NewPool(l)
	for i := 1; i <= 8; i++ {
		go pool.Start()
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handlePublishPostRequest(pool, w, r)
		default:
			w.WriteHeader(403)
		}
	})
}

func main() {
	logger := logging.NewLogger()
	go logger.Start()
	setupRoutes(logger)
	logger.Info("Listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
