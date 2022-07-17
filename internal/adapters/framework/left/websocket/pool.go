package websocket

import (
	"sync"
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/adapters/framework/right/message_queue"
	"github.com/josh-tracey/eventual-agent/internal/logging"
	"github.com/josh-tracey/eventual-agent/internal/profile"
)

var (
	timer = time.NewTicker(120 * time.Second)
)

// Pool - Shared worker pool resources
type Pool struct {
	Subscribe      chan core.SubscribeRequest[*Client]
	Unsubscribe    chan core.SubscribeRequest[*Client]
	UnsubscribeAll chan core.SubscribeRequest[*Client]
	History        chan core.HistoryRequest[*Client]
	Publish        chan core.PublishRequest[*Client]
	core           *core.Adapter
	queue          *message_queue.Adapter
	clientsMap     map[string]*Client
	Logging        *logging.Logger
	cLock          sync.RWMutex
}

// NewPool - Creates new instance of Pool
func NewPool(c *core.Adapter, q *message_queue.Adapter) *Pool {
	return &Pool{
		Subscribe:      make(chan core.SubscribeRequest[*Client], 4),
		Unsubscribe:    make(chan core.SubscribeRequest[*Client], 4),
		UnsubscribeAll: make(chan core.SubscribeRequest[*Client], 4),
		History:        make(chan core.HistoryRequest[*Client], 4),
		Publish:        make(chan core.PublishRequest[*Client], 4),
		core:           c,
		queue:          q,
		clientsMap:     make(map[string]*Client),
		Logging:        c.GetLogger(),
	}
}

func (p *Pool) getClient(refId string) *Client {
	p.cLock.RLock()
	defer p.cLock.RUnlock()
	return p.clientsMap[refId]
}

func (p *Pool) addClient(Id string, c *Client) {
	p.cLock.Lock()
	p.clientsMap[Id] = c
	p.cLock.Unlock()
}

func (p *Pool) Cleaner() {
	for {
		select {
		case <-timer.C:
			p.Logging.Trace("websocket::Pool.Cleaner => Cleaning up clients")
			for _, client := range p.clientsMap {
				if client.closed {
					p.Logging.Trace("websocket::Pool.Cleaner => Removing client %s", client.RefID)
					p.removeClientRefId(client.RefID)
				}
			}
		}
	}
}

func (p *Pool) removeClientRefId(refId string) {
	defer func() {
		if err := recover(); err != nil {
			p.Logging.Error("websocket::Pool.RemoveClientRefId => unhandled exception: %+v", err)
		}
		p.cLock.Unlock()
	}()
	p.cLock.Lock()
	if p.clientsMap[refId] != nil {
		delete(p.clientsMap, refId)
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
			for _, channel := range r.Channels {
				p.queue.Enqueue(channel, r.Event)
				for subChan, sub := range p.core.Subs {
					if channel == subChan || subChan == "global" {
						for _, client := range sub.GetClients() {
							c := p.getClient(client)
							if c != nil && !c.closed {
								p.Logging.Trace("websocket::Pool.Start.Publish => Publishing event to client %s, subscribed to channel %s", client, subChan)
								c.Send <- r.Event
							} else {
								p.Logging.Trace("websocket::Pool.Start.Publish => Client %s is not connected, removing from subscription", client)
								p.removeClientRefId(client)
							}
						}
					}
				}
			}
			profile.Duration(*p.Logging, start, "Pool::Start::Publish")

		case r := <-p.Subscribe:
			p.Logging.Trace("websocket::Pool.Start.Subscribe => Received subscribe event for channels '%s'", r.Channels)
			id := p.core.AddClient(r.Client.ID, r.Channels)
			r.Client.AddRefID(id)
			p.addClient(id, r.Client)
			p.Logging.Trace("websocket::Pool.Start.Subscribe => Added client %s to subscriptions", id)
			p.Logging.Trace("Client Map: %v", p.clientsMap)

		case r := <-p.Unsubscribe:
			p.Logging.Trace("websocket::Pool.Start.Unsubscribe => Received unsubscribe event for channels '%s'", r.Channels)
			p.removeClientRefId(r.Client.RefID)
			p.core.RemoveClientId(r.Client.ID)

		case r := <-p.UnsubscribeAll:
			p.Logging.Trace("websocket::Pool.Start.UnsubscribeAll => Received UnsubscribeAll for %s", r.Client.ID)
			p.removeClientRefId(r.Client.RefID)
			p.core.RemoveClientId(r.Client.ID)

		case r := <-p.History:
			p.Logging.Trace("websocket::Pool.Start.History => Received history event for channel '%s'", r.Channel)
			iter, err := p.queue.Iter(r.Channel, false)
			if err != nil {
				p.Logging.Error("websocket::Pool.Start.History => Error getting iterator for channel '%s': %s", r.Channel, err.Error())
				continue
			}
			for event := range iter {
				r.Client.Send <- event
			}
		}

	}

}
