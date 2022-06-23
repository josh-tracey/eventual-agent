package websocket

import (
	"fmt"
	"net/http"

	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/ports"
)

type Adapter struct {
	core ports.SubjectPort
}

func NewAdapter(core ports.SubjectPort) *Adapter {
	return &Adapter{
		core: core,
	}
}

func (a *Adapter) ListenAndServe(logger *logging.Logger) {
	setupRoutes(logger)
	logger.Info("Listening on 0.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v", err)
		pool.Logging.Warn("websocket upgrade failed: %+v", (err.Error()))
	}

	client := NewClient(r.RemoteAddr, ws, pool)
	pool.Logging.Trace("Received Connection from %+v", client)
	go client.WriteListen()
	client.ReadListen()
}

func setupRoutes(l *logging.Logger) {
	pool := NewPool(l)
	//go pool.CacheManager()
	for i := 1; i <= 32; i++ {
		go pool.Start()
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}
