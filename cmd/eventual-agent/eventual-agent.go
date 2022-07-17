package main

import (
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/left/websocket"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/right/message_queue"
	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/ports"
	"github.com/patrickmn/go-cache"
)

func main() {

	var subs ports.SubjectPort
	var ws ports.WebSocketPort
	var queue ports.MessageQueuePort
	var logger *logging.Logger = logging.NewLogger()

	queue = message_queue.NewAdapter(cache.New(5*time.Minute, 10*time.Minute))

	subs = core.NewAdapter(logger)
	ws = websocket.NewAdapter(subs, queue)

	go logger.Start()
	ws.ListenAndServe()
}
