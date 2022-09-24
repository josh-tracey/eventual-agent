package main

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/left/grpc"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/left/websocket"
	"github.com/josh-tracey/eventual-agent/internal/adapters/services"
	"github.com/josh-tracey/eventual-agent/internal/ports"
	"github.com/josh-tracey/scribe"
)

func main() {
	var subs ports.SubjectPort
	var logger *scribe.Logger = scribe.NewLogger()
	var publishChannel chan *core.PeerEvent = make(chan *core.PeerEvent, 32)
	var subsChannel chan *core.PeerRequest = make(chan *core.PeerRequest, 32)
	var eventQueueChan chan *core.CloudEvent = make(chan *core.CloudEvent, 32)
	var publisher ports.Publisher
	var eventQueue ports.EventQueue

	subs = core.NewAdapter(logger)

	eventQueue, err := services.NewEventQueue(
		subs,
		eventQueueChan,
		subsChannel,
		publishChannel,
	)

	if err != nil {
		panic("EventQueue to Websocket service failed to initialize")
	}

	publisher, err2 := services.NewPublisher(subs, logger, publishChannel)

	if err2 != nil {
		panic("Publisher failed to initialize")
	}

	var ws ports.PeerClient
	ws = websocket.NewAdapter(subs, eventQueueChan)
	grpcServer := grpc.New(subs, logger, publishChannel, subsChannel)

	go logger.Start()
	go grpcServer.Run()
	go ws.ListenAndServe()
	go publisher.Run()
	eventQueue.Run()
}
