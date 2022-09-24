package ports

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
)

type MessageQueuePort interface {
	Enqueue(channel string, message core.CloudEvent)
	Dequeue(channel string, consume bool) (core.CloudEvent, error)
	Iter(channel string, consume bool) (chan core.CloudEvent, error)
}
