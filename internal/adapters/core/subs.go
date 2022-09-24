package core

import (
	"sync"

	"github.com/google/uuid"
)

// Sub - Represents a Subscription
type sub struct {
	ID      string
	clients [][]*string
	lock    sync.RWMutex
}

// NewSub - Creates an instance of Sub
func newSub(id string) *sub {
	return &sub{
		ID:      id,
		clients: [][]*string{},
		lock:    sync.RWMutex{},
	}
}

// AddClient - Thread Safe method of adding a client to Sub
func (s *sub) AddClient(ID *string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	id := uuid.NewString()
	s.clients = append(s.clients, []*string{&id, ID})
	return id
}

func (s *sub) GetClient(ID string) []*string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var clients []*string
	for _, client := range s.clients {
		if s.ID == ID {
			clients = append(clients, client[0])
		}
	}
	return clients
}

// GetClients - Thread Safe method of getting clients from Sub
func (s *sub) GetClients() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var clients []string
	for _, client := range s.clients {
		if client != nil {
			clients = append(clients, *client[1])
		}
	}
	return clients
}

// RemoveClientId - Thread Safe method of removing a client from Sub
func (s *sub) RemoveClientId(client string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for i, clientId := range s.clients {
		if *clientId[1] == client {
			s.clients = removeArr(s.clients, i)
		}
	}
}

func removeArr(s [][]*string, i int) [][]*string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func remove(s []*string, i int) []*string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
