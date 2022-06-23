package ports

type SubjectPort interface {
	AddClient(ID string, channel []string) string
	RemoveClientId(client string)
	RemoveClient(ID string, channel []string)
	GetClients() []string
}
