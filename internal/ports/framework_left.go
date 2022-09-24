package ports

type PeerClient interface {
	ListenAndServe()
}
