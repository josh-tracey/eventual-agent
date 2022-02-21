package websocket

import (
	"sync"
)

var (
	lock sync.Mutex
)

// Sub - Represents a Subscription
type Sub struct {
	ID       string
	channels []string
	clients  []*Client
}

// NewSub - Creates an instance of Sub
func NewSub(id string, channels []string, client *Client) *Sub {
	return &Sub{
		ID:       id,
		channels: channels,
		clients:  []*Client{client},
	}
}

// AddClient - Thread Safe method of adding a client to Sub
func (s *Sub) AddClient(client *Client) {
	lock.Lock()
	s.clients = append(s.clients, client)
	lock.Unlock()
}

// AddNewChannels - Thread Safe method of adding new channels to Sub
func (s *Sub) AddNewChannels(channels []string) {
	lock.Lock()
	var newChannels []string
	for _, channel := range channels {
		var found bool
		for _, subbedchannel := range s.channels {
			if channel == subbedchannel {
				found = true
				break
			}
		}
		if !found {
			newChannels = append(newChannels, channel)
		}
	}
	lock.Unlock()
}

// RemoveChannels - Thread Safe method of removing channels from Sub
func (s *Sub) RemoveChannels(channels []string) {
	lock.Lock()
	for index, subchannel := range s.channels {
		for _, channel := range channels {
			if channel == subchannel {
				s.channels = remove(s.channels, index)
			}
		}
	}
	lock.Unlock()
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
