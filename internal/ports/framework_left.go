package ports

type WebSocketPort interface {
	ListenAndServe()
}
