package ports

import (
	"github.com/josh-tracey/eventual-agent/internal/adapters/core"
	"golang.org/x/net/context"
)

type Publisher interface {
	Run()
	Publish(ctx context.Context, message *core.PeerEvent) error
}

type EventQueue interface {
	Run()
}
