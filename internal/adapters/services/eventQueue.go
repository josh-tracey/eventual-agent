package services

import (
	"errors"
	"sync"
	"time"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/josh-tracey/eventual-agent/internal/ports"
)

type EventQueue struct {
	eventQueueChan chan *core.CloudEvent
	subsChannel    chan *core.PeerRequest
	publishChannel chan *core.PeerEvent
	events         []*core.CloudEvent
	lock           *sync.RWMutex
	subs           *core.Adapter
	timer          time.Ticker
}

func NewEventQueue(
	subs ports.SubjectPort,
	eventQueueChan chan *core.CloudEvent,
	subsChannel chan *core.PeerRequest,
	publishChannel chan *core.PeerEvent,
) (*EventQueue, error) {
	value, err := subs.(*core.Adapter)
	if !err {
		return nil, errors.New("Invalid subject port")
	}

	return &EventQueue{
		events:         []*core.CloudEvent{},
		lock:           &sync.RWMutex{},
		subs:           value,
		eventQueueChan: eventQueueChan,
		subsChannel:    subsChannel,
		publishChannel: publishChannel,
		timer:          *time.NewTicker(time.Microsecond * 100),
	}, nil
}

func (q *EventQueue) Subscribe(peerServer string, channel string, ephemeral bool) {
	q.subs.AddPeer(peerServer, channel, ephemeral)
}

func (eq *EventQueue) AddEvent(event *core.CloudEvent) {
	eq.lock.Lock()
	defer eq.lock.Unlock()

	eq.events = append(eq.events, event)
}

func (eq *EventQueue) GetEvents() (*[]*core.CloudEvent, bool) {
	eq.lock.RLock()
	defer eq.lock.RUnlock()
	numEvents := len(eq.events)

	if numEvents > 0 {
		return &eq.events, true
	}
	return nil, false
}

func (eq *EventQueue) ClearEvents() {
	eq.lock.Lock()
	defer eq.lock.Unlock()

	eq.events = make([]*core.CloudEvent, 0)
}

func (eq *EventQueue) GetEventCount() int {
	eq.lock.RLock()
	defer eq.lock.RUnlock()

	return len(eq.events)
}

func (eq *EventQueue) Run() {

	for {
		select {
		case peerRequest := <-eq.subsChannel:
			eq.Subscribe(peerRequest.PeerAddr, peerRequest.Channel, peerRequest.Ephemeral)
		case event := <-eq.eventQueueChan:
			eq.AddEvent(event)
		case <-eq.timer.C:
			events, ok := eq.GetEvents()
			if ok {
				for _, event := range *events {
					for _, peer := range eq.subs.GetPeerServers() {
						eq.publishChannel <- &core.PeerEvent{
							PeerServer: peer.ID,
							Event:      *event,
						}
					}
				}
				eq.ClearEvents()
			}
		}
	}

}
