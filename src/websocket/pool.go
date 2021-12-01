package websocket

import (
	"sync"
	"time"

	"github.com/josh-tracey/eventual-agent/src/logging"
	"github.com/josh-tracey/eventual-agent/src/profile"
)

type Pool struct {
	Subcribe       chan SubscribeRequest
	Unsubscribe    chan SubscribeRequest
	UnsubscribeAll chan SubscribeRequest
	Publish        chan PublishRequest
	Channels       map[string][]*Client
	Logging        *logging.Logger
}

var (
	channelMutex sync.Mutex
)

func NewPool(logger *logging.Logger) *Pool {
	return &Pool{
		Subcribe:       make(chan SubscribeRequest),
		Unsubscribe:    make(chan SubscribeRequest),
		UnsubscribeAll: make(chan SubscribeRequest),
		Publish:        make(chan PublishRequest),
		Channels:       make(map[string][]*Client),
		Logging:        logger,
	}
}

func (p *Pool) channelHasClient(channel string, client *Client) (chan bool, chan int) {

	var found chan bool = make(chan bool)
	var index chan int = make(chan int)

	go func() {
		defer profile.Duration(*p.Logging, time.Now(), "channelHasClient")
		for i, cli := range p.Channels[channel] {
			if cli == client {
				found <- true
				index <- i
				return
			}
		}
		found <- false
		index <- -1
	}()

	return found, index
}

func (p *Pool) removeClientFromChannel(channel string, client *Client) chan []*Client {

	var result chan []*Client

	go func() {
		defer func() {
			channelMutex.Unlock()
			profile.Duration(*p.Logging, time.Now(), "removeClientFromChannel")
		}()
		channelMutex.Lock()
		clients := p.Channels[channel]

		for i, cli := range clients {
			if cli == client {
				result <- remove(clients, i)
				return
			}
		}
		result <- p.Channels[channel]
	}()

	return result
}

func (p *Pool) removeClientFromAllChannels(client *Client) {
	defer func() {
		channelMutex.Unlock()
		profile.Duration(*p.Logging, time.Now(), "removeClientFromAllChannels")
	}()

	channelMutex.Lock()
	for channelName := range p.Channels {

		for i, cli := range p.Channels[channelName] {
			if cli == client {
				p.Channels[channelName] = remove(p.Channels[channelName], i)
			}
		}
	}
}

func remove(s []*Client, i int) []*Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (p *Pool) Start() {
	for {
		select {

		case r := <-p.Publish:
			// This is potentially sluggish if publishing to many channels
			p.Logging.Debug("Received publish event for channel '%s'", r.Event.Source)
			for _, channelName := range r.Channels {
				if channelName == "global" {
					continue
				}
				for i, client := range p.Channels[channelName.(string)] {
					if err := <-client.Send(r.Event); err != nil {
						channelMutex.Lock()
						p.Logging.Error(err.Error())
						p.Channels[channelName.(string)] = remove(p.Channels[channelName.(string)], i)
						channelMutex.Unlock()
					}
				}
			}
			for i, client := range p.Channels["global"] {
				if err := <-client.Send(r.Event); err != nil {
					channelMutex.Lock()
					p.Logging.Error(err.Error())
					p.Channels["global"] = remove(p.Channels["global"], i)
					channelMutex.Unlock()
				}
			}

		case r := <-p.Subcribe:
			channelMutex.Lock()
			for i := range r.Channels {
				func(channel interface{}) {
					found, _ := p.channelHasClient(channel.(string), r.Client)
					b := <-found
					if !b {
						p.Logging.Debug("Adding client %s to channel: %s", r.Client.ID, channel)
						p.Channels[channel.(string)] = append(p.Channels[channel.(string)], r.Client)
					}
				}(r.Channels[i])
				channelMutex.Unlock()
			}

		case r := <-p.Unsubscribe:
			channelMutex.Lock()
			for _, channel := range r.Channels {
				p.Logging.Debug("Unsubscribing client %s from channel: %s", r.Client.ID, channel)
				result := p.removeClientFromChannel(channel.(string), r.Client)
				p.Channels[channel.(string)] = <-result
			}
			channelMutex.Unlock()

		// This is potentially slow, if large amount of channels active in memory. O(n^2)
		case r := <-p.UnsubscribeAll:
			p.Logging.Debug("Unsubscribing client %s from all channels", r.Client.ID)
			p.removeClientFromAllChannels(r.Client)
		}
	}
}
