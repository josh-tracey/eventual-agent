package ports

import (
	"github.com/josh-tracey/eventual-agent/internal/logging"
)

type WebSocketPort interface {
	ListenAndServe(logger *logging.Logger)
}
