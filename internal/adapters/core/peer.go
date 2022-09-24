package core

import "sync"

type peer struct {
	ID       string
	channels []*string
	clients  []*string
	pLock    sync.RWMutex
}

// Newpeer - Creates an instance of peer
func newPeer(id string) *peer {
	return &peer{
		ID:       id,
		channels: []*string{},
		clients:  []*string{},
		pLock:    sync.RWMutex{},
	}
}

func (p *peer) GetChannels() []*string {
	return p.channels
}

// AddClient - Thread Safe method of adding a client to peer
func (p *peer) AddChannel(channel string) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.channels = append(p.channels, &channel)
}

// RemoveChannel - Thread Safe method of removing a client from peer
func (p *peer) RemoveChannel(channel string) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	for i, c := range p.channels {
		if *c == channel {

			p.channels = remove(p.channels, i)

		}
	}

}

func (p *peer) GetClients() []*string {
	return p.clients
}

func (p *peer) AddClient(client string) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	p.clients = append(p.clients, &client)
}

func (p *peer) RemoveClient(client string) {
	p.pLock.Lock()
	defer p.pLock.Unlock()
	for i, c := range p.clients {
		if *c == client {

			p.clients = remove(p.clients, i)

		}
	}

}

func (p *peer) GetChannelCount() int {
	return len(p.channels)
}

func (p *peer) GetClientCount() int {
	return len(p.clients)
}

func (p *peer) IterChannels() chan *string {
	c := make(chan *string)
	go func() {
		p.pLock.RLock()
		defer p.pLock.RUnlock()
		for _, channel := range p.channels {
			c <- channel
		}
		close(c)
	}()
	return c
}

func (p *peer) IterClients() chan *string {
	c := make(chan *string)
	go func() {
		p.pLock.RLock()
		defer p.pLock.RUnlock()
		for _, client := range p.clients {
			c <- client
		}
		close(c)
	}()
	return c
}
