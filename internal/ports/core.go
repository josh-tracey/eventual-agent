package ports

import "github.com/josh-tracey/scribe"

type SubjectPort interface {
	AddPeer(addr string, channel string, ephemeral bool) string
	RemoveClientId(client string)
	RemoveClient(ID string, channel []string)
	GetClients() []string
	GetLogger() *scribe.Logger
}
