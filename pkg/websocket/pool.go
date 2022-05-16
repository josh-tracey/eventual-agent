package websocket

import (
	"sync"
	"time"

	"github.com/josh-tracey/eventual-agent/pkg/logging"
	"github.com/josh-tracey/eventual-agent/pkg/profile"
)

// TODO: Cache Manager
var (
	sublock  sync.Mutex
	snapshot map[string]*Sub = make(map[string]*Sub)
	channels map[string]*Client
	timer    = time.NewTicker(120 * time.Second)
)

// Pool - Shared worker pool resources
type Pool struct {
	Subscribe      chan SubscribeRequest
	Unsubscribe    chan SubscribeRequest
	UnsubscribeAll chan SubscribeRequest
	Publish        chan PublishRequest
	Subs           map[string]*Sub
	Logging        *logging.Logger
}

// NewPool - Creates new instance of Pool
func NewPool(logger *logging.Logger) *Pool {
	return &Pool{
		Subscribe:      make(chan SubscribeRequest),
		Unsubscribe:    make(chan SubscribeRequest),
		UnsubscribeAll: make(chan SubscribeRequest),
		Publish:        make(chan PublishRequest),
		Subs:           make(map[string]*Sub),
		Logging:        logger,
	}
}

// TakeSnapshot - Takes Snapshot of Current State of Subscriptions. TODO: CacheManager
func (p *Pool) TakeSnapshot() {
	//sublock.lock()
	//for _, _ := range p.subs {
	////
	//}
	//sublock.unlock()
}

//CacheManager - Managers servers cache TODO: CacheManager
func (p *Pool) CacheManager() {

	p.Logging.Info("Cache Manager Started...")

	for {
		select {
		case <-timer.C:
			p.Logging.Trace("CacheManager::takesnapshot")
			p.TakeSnapshot()
			p.Logging.Debug("CacheManager::snapshot: %v", snapshot)
		}
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
			p.Logging.Trace("Received publish event for channel '%s'", r.Event.Type)
			//TODO: Replace with more efficient publishing method... O(x^4) TODO: CacheManager
			for _, sub := range p.Subs {
				for _, channel := range r.Channels {
					for _, subChan := range sub.channels {
						if channel == subChan || subChan == "global" {
							for _, client := range sub.clients {
								if !client.closed {
									p.Logging.Trace("SENDING... %v", client.ID)
									client.Send <- r.Event
								}
							}
						}
					}
				}
			}
			profile.Duration(*p.Logging, start, "Pool::Start::Publish")

		case r := <-p.Subscribe:
			sublock.Lock()
			p.Logging.Info("Received subscribe event for channels '%s'", r.Channels)
			if sub, ok := p.Subs[r.Client.ID]; ok { //O(x^2)
				//	check / update
				sub.AddNewChannels(r.Channels) // O(x^2)
				sub.AddClient(r.Client)        // O(1)
			} else { // O(1)
				//	create sub
				sub := NewSub(r.Client.ID, r.Channels, r.Client)
				p.Subs[r.Client.ID] = sub
			}
			p.Logging.Trace("Subs '%v'", p.Subs)
			sublock.Unlock()

		case r := <-p.Unsubscribe:
			sublock.Lock()
			p.Logging.Trace("Received unsubscribe event for channels '%s'", r.Channels)

			if sub, ok := p.Subs[r.Client.ID]; ok { //O(x^2)
				sub.RemoveChannels(r.Channels) // O(x^2)
				if len(sub.channels) == 0 {
					delete(p.Subs, r.Client.ID)
				}
			}
			sublock.Unlock()

		case r := <-p.UnsubscribeAll:
			sublock.Lock()
			p.Logging.Trace("Received UnsubscribeAll for %s", r.Client.ID)
			delete(p.Subs, r.Client.ID) // O(1)
			sublock.Unlock()
		}
	}
}
