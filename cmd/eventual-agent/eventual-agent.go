package main

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/left/websocket"
	"github.com/josh-tracey/eventual-agent/internal/ports"
	"github.com/josh-tracey/notary"
	"github.com/josh-tracey/scribe"
)

func main() {
	n := notary.New("secret")
	var subs ports.SubjectPort
	var ws ports.WebSocketPort
	var logger *scribe.Logger = scribe.NewLogger()

	subs = core.NewAdapter(logger)
	ws = websocket.NewAdapter(subs)

	token, err := n.NewSignedToken()

	go logger.Start()
	valid, err := n.VerifyToken(token)
	if err != nil {
		logger.Error("%v", err)
	}
	if !valid {
		panic("Invalid token")
	}
	ws.ListenAndServe()
}
