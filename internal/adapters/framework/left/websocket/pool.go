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
	Subs           *core.Adapter
	clientsMap     map[string]*Client
	Logging        *logging.Logger
}

// NewPool - Creates new instance of Pool
func NewPool(logger *logging.Logger) *Pool {
	return &Pool{
		Subscribe:      make(chan SubscribeRequest),
		Unsubscribe:    make(chan SubscribeRequest),
		UnsubscribeAll: make(chan SubscribeRequest),
		Publish:        make(chan PublishRequest),
		Subs:           core.NewAdapter(),
		clientsMap:     make(map[string]*Client),
		Logging:        logger,
	}
}

// Start - Go Routine runs worker with shared Pool resources.
func (p *Pool) Start() {

	defer func() {
		if err := recover(); err != nil {
			p.Logging.Error("unhandled exception in pool: %+v", err)
		}
		p.Logging.Warn("Worker stopped")
		p.Start()
	}()

	for {
		select {

		case r := <-p.Publish:
			start := time.Now()
			p.Logging.Trace("Received publish event for channel '%s'", r.Event)
			for _, channel := range r.Channels {
				for subChan, sub := range p.Subs.Subs {
					if channel == subChan || subChan == "global" {
						for _, client := range sub.GetClients() {
							c := p.clientsMap[client]
							if !c.closed {
								p.Logging.Debug("Publishing event to client %s, subscribed to channel %s", client, subChan)
								c.Send <- r.Event
							}
						}
					}
				}
			}
			profile.Duration(*p.Logging, start, "Pool::Start::Publish")

		case r := <-p.Subscribe:
			p.Logging.Info("Received subscribe event for channels '%s'", r.Channels)
			id := p.Subs.AddClient(r.Client.ID, r.Channels)
			p.clientsMap[id] = r.Client
			p.Logging.Info("Added client %s to subscriptions", id)
			p.Logging.Info("Client Map: %v", p.clientsMap)

		case r := <-p.Unsubscribe:
			p.Logging.Trace("Received unsubscribe event for channels '%s'", r.Channels)
			p.Subs.RemoveClientId(r.Client.ID)

		case r := <-p.UnsubscribeAll:
			p.Logging.Trace("Received UnsubscribeAll for %s", r.Client.ID)
			p.Subs.RemoveClientId(r.Client.ID)
		}
	}
}
