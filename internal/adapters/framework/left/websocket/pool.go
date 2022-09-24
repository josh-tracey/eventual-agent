package websocket

import (
	"sync"
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/scribe"
)

var (
	timer = time.NewTicker(120 * time.Second)
)

// Pool - Shared worker pool resources
type Pool struct {
	Subscribe      chan core.SubscribeRequest[*Client]
	Unsubscribe    chan core.SubscribeRequest[*Client]
	UnsubscribeAll chan core.SubscribeRequest[*Client]
	Publish        chan core.PublishRequest[*Client]
	core           *core.Adapter
	clientsMap     *sync.Map
	Logging        *scribe.Logger
	cLock          *sync.RWMutex
	grpcEventQueue chan *core.CloudEvent
}

// NewPool - Creates new instance of Pool
func NewPool(c *core.Adapter, grpcEventQueue chan *core.CloudEvent) *Pool {
	return &Pool{
		Subscribe:      make(chan core.SubscribeRequest[*Client], 4),
		Unsubscribe:    make(chan core.SubscribeRequest[*Client], 4),
		UnsubscribeAll: make(chan core.SubscribeRequest[*Client], 4),
		Publish:        make(chan core.PublishRequest[*Client], 4),
		core:           c,
		clientsMap:     &sync.Map{},
		Logging:        c.GetLogger(),
		cLock:          &sync.RWMutex{},
		grpcEventQueue: grpcEventQueue,
	}
}

func (p *Pool) getClient(refId string) *Client {
	p.cLock.RLock()
	defer p.cLock.RUnlock()
	client, found := p.clientsMap.Load(refId)
	if !found {
		return nil
	}

	return client.(*Client)
}

func (p *Pool) addClient(Id string, c *Client) {
	p.clientsMap.Store(Id, c)
}

func (p *Pool) Cleaner() {
	for {
		select {
		case <-timer.C:
			p.Logging.Trace("websocket::Pool.Cleaner => Cleaning up clients")
			p.clientsMap.Range(func(id, client interface{}) bool {
				if client.(*Client).closed {
					p.Logging.Trace("websocket::Pool.Cleaner => Removing client %s", client.(*Client).RefID)
					p.removeClientRefId(client.(*Client).RefID)
				}
				return true
			})
		}
	}
}

func (p *Pool) removeClientRefId(refId string) {
	defer func() {
		if err := recover(); err != nil {
			p.Logging.Error("websocket::Pool.RemoveClientRefId => unhandled exception: %+v", err)
		}
	}()
	p.clientsMap.Delete(refId)
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

			go func() {
				p.grpcEventQueue <- &r.Event
			}()

			for _, sub := range p.core.GetPeerClients() {
				for _, channel := range sub.GetChannels() {
					if *channel == r.Event.Type || *channel == "global" {
						subscription := p.core.GetSub(*channel)
						if subscription != nil {
							p.Logging.Trace("websocket::Pool.Start.Publish => Sending event to client %v", sub.ID)
							client := subscription.GetClient(*channel)
							for i, c := range client {
								c := p.getClient(*c)
								if c != nil && !c.closed {
									p.Logging.Trace("websocket::Pool.Start.Publish => Publishing event to client %v, subscribed to channel %v", *client[i], channel)
									c.Send <- r.Event
								} else {
									p.Logging.Trace("websocket::Pool.Start.Publish => Client %s is not connected, removing from subscription", *client[i])
									p.removeClientRefId(*client[i])
								}
							}
						}
					}
				}
			}

			p.Logging.Duration(start, "Pool::Start::Publish")

		case r := <-p.Subscribe:
			p.Logging.Trace("websocket::Pool.Start.Subscribe => Received subscribe event for channels '%s'", r.Channels)
			for _, channel := range r.Channels {
				refID := p.core.AddPeer(r.Client.ID, channel, true)
				p.addClient(refID, r.Client)
				p.Logging.Trace("websocket::Pool.Start.Subscribe => Added client %v to subscriptions", p.core.GetPeerClients())
			}

		case r := <-p.Unsubscribe:
			p.Logging.Trace("websocket::Pool.Start.Unsubscribe => Received unsubscribe event for channels '%s'", r.Channels)
			p.removeClientRefId(r.Client.RefID)
			for _, channel := range r.Channels {
				p.core.RemovePeer(channel, true)
			}

		case r := <-p.UnsubscribeAll:
			p.Logging.Trace("websocket::Pool.Start.UnsubscribeAll => Received UnsubscribeAll for %s", r.Client.ID)
			p.removeClientRefId(r.Client.RefID)
			for _, channel := range r.Channels {
				p.core.RemovePeer(channel, true)
			}
		}
	}
}
