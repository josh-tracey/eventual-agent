package core

import (
	"sync"

	"github.com/josh-tracey/scribe"
)

type Adapter struct {
	logger      *scribe.Logger
	subs        map[string]*sub
	peerClients map[string]*peer
	peerServers map[string]*peer
	lock        sync.RWMutex
}

func NewAdapter(logger *scribe.Logger) *Adapter {
	return &Adapter{
		logger:      logger,
		subs:        map[string]*sub{"global": newSub("global")},
		peerServers: make(map[string]*peer),
		peerClients: make(map[string]*peer),
		lock:        sync.RWMutex{},
	}
}

func (adapt *Adapter) GetLogger() *scribe.Logger {
	return adapt.logger
}

// GetSub - Thread Safe method of getting a Sub
func (adapt *Adapter) GetSub(channel string) *sub {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.GetSub => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()
	if _, ok := adapt.subs[channel]; !ok {
		adapt.subs[channel] = newSub(channel)
	}

	return adapt.subs[channel]
}

func (adapt *Adapter) GetPeer(addr string, ephemeral bool) *peer {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.GetPeer => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()

	if !ephemeral {

		if _, ok := adapt.peerServers[addr]; !ok {
			adapt.peerServers[addr] = newPeer(addr)
		}

		return adapt.peerServers[addr]
	} else {
		if _, ok := adapt.peerClients[addr]; !ok {
			adapt.peerClients[addr] = newPeer(addr)
		}

		return adapt.peerClients[addr]
	}
}

func (adapt *Adapter) RemovePeer(addr string, ephemeral bool) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemovePeer => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()

	if !ephemeral {
		delete(adapt.peerServers, addr)
	} else {
		delete(adapt.peerClients, addr)
	}
}

func (adapt *Adapter) RemovePeerFromChannel(addr string, channel string, ephemeral bool) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemovePeerFromChannel => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()

	if !ephemeral {
		adapt.peerServers[addr].RemoveChannel(channel)
	} else {
		adapt.peerClients[addr].RemoveChannel(channel)
	}
}

func (adapt *Adapter) RemoveClientFromChannel(client string, channel string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemoveClientFromChannel => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()

	adapt.subs[channel].RemoveClientId(client)
}

func (adapt *Adapter) HasPeerId(addr string, ephemeral bool) bool {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.hasPeerId => %s", r)
		}
	}()

	adapt.lock.RLock()
	defer adapt.lock.RUnlock()

	if !ephemeral {
		if _, ok := adapt.peerServers[addr]; ok {
			return true
		}
	} else {
		if _, ok := adapt.peerClients[addr]; ok {
			return true
		}
	}

	return false
}

func (adapt *Adapter) GetPeerServers() map[string]*peer {
	adapt.lock.RLock()
	defer adapt.lock.RUnlock()
	return adapt.peerServers
}

func (adapt *Adapter) GetPeerClients() map[string]*peer {
	adapt.lock.RLock()
	defer adapt.lock.RUnlock()
	return adapt.peerClients
}

func (adapt *Adapter) GetSubs() map[string]*sub {
	adapt.lock.RLock()
	defer adapt.lock.RUnlock()
	return adapt.subs
}

// Remove - Thread Safe method of removing a channel from Sub
func (adapt *Adapter) Remove(channel string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.Remove => %s", r)
		}
	}()

	adapt.lock.Lock()
	defer adapt.lock.Unlock()
	delete(adapt.subs, channel)
}

func (adapt *Adapter) RemoveClientId(client string) {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.RemoveClientId => %s", r)
		}
	}()

	var done chan bool = make(chan bool)
	go func() {
		adapt.lock.Lock()
		defer adapt.lock.Unlock()
		for _, v := range adapt.subs {
			for i, clientId := range v.clients {
				if *clientId[1] == client {
					v.clients = removeArr(v.clients, i)
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
	for _, sub := range adapt.subs {
		clients = append(clients, sub.GetClients()...)
	}
	return clients
}

func (adapt *Adapter) AddPeer(addr string, channel string, ephemeral bool) string {
	defer func() {
		if r := recover(); r != nil {
			adapt.logger.Error("core::Adapter.AddPeer => %s", r)
		}
	}()

	peer := adapt.GetPeer(addr, ephemeral)
	peer.AddChannel(channel)
	if ephemeral {
		if adapt.subs[channel] == nil {
			adapt.subs[channel] = newSub(channel)
		}
		return adapt.subs[channel].AddClient(&addr)
	}
	return ""
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
