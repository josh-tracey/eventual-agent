package main

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/left/websocket"
	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/ports"
)

func main() {

	var subs ports.SubjectPort
	var ws ports.WebSocketPort
	subs = core.NewAdapter()
	ws = websocket.NewAdapter(subs)

	logger := logging.NewLogger()
	go logger.Start()
	ws.ListenAndServe(logger)
}
