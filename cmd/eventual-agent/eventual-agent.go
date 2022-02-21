package main

import (
	"fmt"
	"net/http"

	"github.com/josh-tracey/eventual-agent/pkg/logging"
	"github.com/josh-tracey/eventual-agent/pkg/websocket"
)

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		pool.Logging.Warn("websocket upgrade failed: %+v", (err.Error()))
	}

	client := websocket.NewClient(r.RemoteAddr, ws, pool)
	pool.Logging.Trace("Received Connection from %+v", client)
	go client.WriteListen()
	client.ReadListen()
}

func setupRoutes(l *logging.Logger) {
	pool := websocket.NewPool(l)
	//go pool.CacheManager()
	for i := 1; i <= 32; i++ {
		go pool.Start()
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	logger := logging.NewLogger()
	go logger.Start()
	setupRoutes(logger)
	logger.Info("Listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
