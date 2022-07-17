package message_queue

// Using a cache to store the messages in the queue, for simplicity.
// https://github.com/patrickmn/go-cache

import (
	"errors"
	"fmt"
	"log"

	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"github.com/patrickmn/go-cache"
)

var ()

type Adapter struct {
	cache *cache.Cache
}

func NewAdapter(c *cache.Cache) *Adapter {
	adapter := &Adapter{}
	adapter.init(c)
	return adapter
}

func (adapter *Adapter) init(c *cache.Cache) {
	adapter.cache = c
}

func (adapter *Adapter) Enqueue(channel string, message core.CloudEvent) {
	eventQueue, found := adapter.cache.Get(channel)
	if !found {
		log.Printf("Creating new queue for channel %s", channel)
		eventQueue = []core.CloudEvent{message}
	}
	q, _ := eventQueue.([]core.CloudEvent)
	q = append(eventQueue.([]core.CloudEvent), message)

	adapter.cache.Set(channel, q, cache.DefaultExpiration)
}

func (adapter *Adapter) Dequeue(channel string, consume bool) (core.CloudEvent, error) {
	eventQueue, found := adapter.cache.Get(channel)
	if !found {
		return core.CloudEvent{}, errors.New(fmt.Sprintf("Channel %s queue not found", channel))
	}

	q, _ := eventQueue.([]core.CloudEvent)
	if len(q) == 0 {
		return core.CloudEvent{}, errors.New(fmt.Sprintf("CloudEvents for channel %s not found", channel))
	}

	message := q[0]
	if consume {
		q = q[1:]
	}

	adapter.cache.Set(channel, q, cache.DefaultExpiration)
	return message, nil
}

func (adapter *Adapter) Iter(channel string, consume bool) (chan core.CloudEvent, error) {

	eventQueue, found := adapter.cache.Get(channel)
	if !found {
		return nil, errors.New(fmt.Sprintf("Channel %s queue not found", channel))
	}

	q, _ := eventQueue.([]core.CloudEvent)
	if len(q) == 0 {
		return nil, errors.New(fmt.Sprintf("CloudEvents for channel %s not found", channel))
	}

	if consume {
		adapter.cache.Delete(channel)
	}

	c := make(chan core.CloudEvent)
	go func() {
		for _, message := range q {
			c <- message
		}
		close(c)
	}()
	return c, nil
}
