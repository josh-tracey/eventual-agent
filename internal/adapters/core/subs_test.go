package core

import (
	"testing"

	"github.com/josh-tracey/eventual-agent/internal/logging"
)

func TestSubs(t *testing.T) {
	logger := logging.NewLogger()
	subs := NewAdapter(logger)
	subs.AddClient("12345", []string{"tempUpdates", "buttonPresses"})
	subs.GetClients()
}
