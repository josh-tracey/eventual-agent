package core

import (
	"sync"

	"github.com/google/uuid"
	"github.com/josh-tracey/eventual-agent/internal/logging"
)

var (
	lock sync.RWMutex
)

type Adapter struct {
	logger *logging.Logger
	Subs   map[string]*sub
}

// Sub - Represents a Subscription
type sub struct {
	ID      string
	clients []*string
}

func NewAdapter(logger *logging.Logger) *Adapter {
	return &Adapter{
		logger: logger,
		Subs:   make(map[string]*sub),
	}
}

func (adapt *Adapter) GetLogger() *logging.Logger {
	return adapt.logger
}

// GetSub - Thread Safe method of getting a Sub
func (adapt *Adapter) GetSub(channel string) *sub {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.GetSub => %s", r)
		}
		lock.RUnlock()
	}()

	lock.RLock()
	if _, ok := adapt.Subs[channel]; !ok {
		adapt.Subs[channel] = newSub(channel, nil)
	}

	return adapt.Subs[channel]
}

// Remove - Thread Safe method of removing a channel from Sub
func (adapt *Adapter) Remove(channel string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.Remove => %s", r)
		}
		lock.Unlock()
	}()

	lock.Lock()
	delete(adapt.Subs, channel)
}

func (adapt *Adapter) RemoveClientId(client string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemoveClientId => %s", r)
		}
	}()

	var done chan bool = make(chan bool)
	go func() {
		for _, v := range adapt.Subs {
			for _, clientId := range v.clients {
				if clientId != nil {
					v.RemoveClientId(*clientId)
				}
			}
		}
		done <- true
		close(done)
	}()
	<-done
}

func (adapt *Adapter) GetClients() []string {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.GetClients => %s", r)
		}
	}()

	var clients []string
	for _, sub := range adapt.Subs {
		clients = append(clients, sub.GetClients()...)
	}
	return clients
}

func (adapt *Adapter) AddClient(ID string, channels []string) string {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.AddClients => %s", r)
		}
	}()

	var client string
	for _, channel := range channels {
		sub := adapt.GetSub(channel)
		client = sub.AddClient(&ID)
	}
	return client
}

func (adapt *Adapter) RemoveClient(ID string, channels []string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemoveClient => %s", r)
		}
	}()

	for _, channel := range channels {
		adapt.GetSub(channel).RemoveClientId(ID)
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

// GetClients - Thread Safe method of getting clients from Sub
func (s *sub) GetClients() []string {
	lock.RLock()
	var clients []string
	for _, client := range s.clients {
		if client != nil {
			clients = append(clients, *client)
		}
	}
	lock.RUnlock()
	return clients
}

// RemoveClientId - Thread Safe method of removing a client from Sub
func (s *sub) RemoveClientId(client string) {
	lock.Lock()
	for i, clientId := range s.clients {
		if clientId != nil && *clientId == client {
			remove(s.clients, i)
		}
	}
	lock.Unlock()
}

func remove(s []*string, i int) []*string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
