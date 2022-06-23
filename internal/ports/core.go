package ports

import "github.com/josh-tracey/eventual-agent/internal/logging"

type SubjectPort interface {
	AddClient(ID string, channel []string) string
	RemoveClientId(client string)
	RemoveClient(ID string, channel []string)
	GetClients() []string
	GetLogger() *logging.Logger
}
