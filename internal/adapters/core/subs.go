package core

import (
	"sync"

	"github.com/google/uuid"
)

var (
	lock sync.Mutex
)

type Adapter struct {
	Subs map[string]*sub
}

// Sub - Represents a Subscription
type sub struct {
	ID      string
	clients []*string
}

func NewAdapter() *Adapter {
	return &Adapter{
		Subs: make(map[string]*sub),
	}
}

// GetSub - Thread Safe method of getting a Sub
func (adapt *Adapter) GetSub(channel string) *sub {
	lock.Lock()
	if _, ok := adapt.Subs[channel]; !ok {
		adapt.Subs[channel] = newSub(channel, nil)
	}
	lock.Unlock()
	return adapt.Subs[channel]
}

// Remove - Thread Safe method of removing a channel from Sub
func (adapt *Adapter) Remove(channel string) {
	lock.Lock()
	delete(adapt.Subs, channel)
	lock.Unlock()
}

func (adapt *Adapter) RemoveClientId(client string) {
	lock.Lock()
	var done chan bool = make(chan bool)
	go func() {
		for _, v := range adapt.Subs {
			for i, clientId := range v.clients {
				if clientId != nil && *clientId == client {
					remove(v.clients, i)
					done <- true
				}
			}
		}
		close(done)
	}()
	<-done
	lock.Unlock()
}

func (adapt *Adapter) GetClients() []string {
	var clients []string
	for _, sub := range adapt.Subs {
		clients = append(clients, sub.GetClients()...)
	}
	return clients
}

func (adapt *Adapter) AddClient(ID string, channels []string) string {
	var client string
	for _, channel := range channels {
		sub := adapt.GetSub(channel)
		client = sub.AddClient(&ID)
	}
	return client
}

func (adapt *Adapter) RemoveClient(ID string, channels []string) {
	for _, channel := range channels {
		adapt.GetSub(channel).RemoveClient(&ID)
	}
}

// NewSub - Creates an instance of Sub
func newSub(id string, client *string) *sub {
	return &sub{
		ID:      id,
		clients: []*string{client},
	}
}

// AddClient - Thread Safe method of adding a client to Sub
func (s *sub) AddClient(ID *string) string {
	lock.Lock()
	id := uuid.New().String()
	s.clients = append(s.clients, &id)
	lock.Unlock()
	return id
}

// RemoveClient - Thread Safe method of removing a client from Sub
func (s *sub) RemoveClient(ID *string) {
	lock.Lock()
	for index, subclient := range s.clients {
		if *subclient == *ID {
			s.clients = remove(s.clients, index)
		}
	}
	lock.Unlock()
}

// GetClients - Thread Safe method of getting clients from Sub
func (s *sub) GetClients() []string {
	lock.Lock()
	var clients []string
	for _, client := range s.clients {
		if client != nil {
			clients = append(clients, *client)
		}
	}
	lock.Unlock()
	return clients
}

// RemoveClientId - Thread Safe method of removing a client from Sub
func (s *sub) RemoveClientId(client string) {
	lock.Lock()
	for index, subclient := range s.clients {
		if *subclient == client {
			s.clients = remove(s.clients, index)
		}
	}
	lock.Unlock()
}

func remove(s []*string, i int) []*string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
