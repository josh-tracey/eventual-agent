package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/right/message_queue"
	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/ports"
)

type Adapter struct {
	core  *core.Adapter
	queue *message_queue.Adapter
}

func NewAdapter(c ports.SubjectPort, q ports.MessageQueuePort) *Adapter {
	value, ok := c.(*core.Adapter)
	if !ok {
		c.GetLogger().Error("websocket::Adapter.NewAdapter => Failed to cast c to *core.Adapter")
	}

	qValue, ok := q.(*message_queue.Adapter)
	if !ok {
		c.GetLogger().Error("websocket::Adapter.NewAdapter => Failed to cast q to *message_queue.Adapter")
	}

	return &Adapter{
		core:  value,
		queue: qValue,
	}
}

func (a *Adapter) ListenAndServe() {
	setupRoutes(a)
	a.core.GetLogger().Info("Listening on 0.0.0.0:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(logging.FgRed, "Fatal: ", logging.Reset, err)
	}
}

func serveWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			pool.Logging.Error("websocket::Pool.serveWs => %s", r)
		}
	}()

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

func setupRoutes(a *Adapter) {
	pool := NewPool(a.core, a.queue)
	for i := 1; i <= 32; i++ {
		go pool.Start()
	}
	go pool.Cleaner()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}
