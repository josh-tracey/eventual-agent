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
	var logger *logging.Logger

	logger = logging.NewLogger()
	subs = core.NewAdapter(logger)
	ws = websocket.NewAdapter(subs)

	go logger.Start()
	ws.ListenAndServe()
}
