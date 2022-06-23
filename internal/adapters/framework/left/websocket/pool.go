package websocket

import (
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/profile"
)

// Pool - Shared worker pool resources
type Pool struct {
	Subscribe      chan SubscribeRequest
	Unsubscribe    chan SubscribeRequest
	UnsubscribeAll chan SubscribeRequest
	Publish        chan PublishRequest
	core           *core.Adapter
	clientsMap     map[string]*Client
	Logging        *logging.Logger
}

// NewPool - Creates new instance of Pool
func NewPool(c *core.Adapter) *Pool {
	return &Pool{
		Subscribe:      make(chan SubscribeRequest, 4),
		Unsubscribe:    make(chan SubscribeRequest, 4),
		UnsubscribeAll: make(chan SubscribeRequest, 4),
		Publish:        make(chan PublishRequest, 4),
		core:           c,
		clientsMap:     make(map[string]*Client),
		Logging:        c.GetLogger(),
	}
}

// Start - Go Routine runs worker with shared Pool resources.
func (p *Pool) Start() {

	defer func() {
		if err := recover(); err != nil {
			p.Logging.Error("websocket::Pool.Start => unhandled exception: %+v", err)
		}
		p.Logging.Warn("Worker stopped")
		p.Start()
	}()

	for {
		select {

		case r := <-p.Publish:
			start := time.Now()
			p.Logging.Trace("websocket::Pool.Start.Publish => Received publish event for channel '%s'", r.Event.Type)
			for subChan, sub := range p.core.Subs {
				for _, channel := range r.Channels {
					if channel == subChan || subChan == "global" {
						for _, client := range sub.GetClients() {
							c := p.clientsMap[client]
							if !c.closed {
								p.Logging.Debug("websocket::Pool.Start.Publish => Publishing event to client %s, subscribed to channel %s", client, subChan)
								c.Send <- r.Event
							}
						}
					}
				}
			}
			profile.Duration(*p.Logging, start, "Pool::Start::Publish")

		case r := <-p.Subscribe:
			p.Logging.Info("websocket::Pool.Start.Subscribe => Received subscribe event for channels '%s'", r.Channels)
			id := p.core.AddClient(r.Client.ID, r.Channels)
			p.clientsMap[id] = r.Client
			p.Logging.Info("websocket::Pool.Start.Subscribe => Added client %s to subscriptions", id)
			p.Logging.Info("Client Map: %v", p.clientsMap)

		case r := <-p.Unsubscribe:
			p.Logging.Trace("websocket::Pool.Start.Unsubscribe => Received unsubscribe event for channels '%s'", r.Channels)
			p.core.RemoveClientId(r.Client.ID)

		case r := <-p.UnsubscribeAll:
			p.Logging.Trace("websocket::Pool.Start.UnsubscribeAll => Received UnsubscribeAll for %s", r.Client.ID)
			p.core.RemoveClientId(r.Client.ID)
		}
	}
}
