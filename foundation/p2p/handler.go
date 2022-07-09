package p2p

// Handler handles requests to the p2p network.
type Handler interface {
	Serve(ResponseWriter, *Request) error
}

// HandlerFunc wraps a Handler into a function.
type HandlerFunc func(ResponseWriter, *Request) error

// Serve implements the Handler interface for HandlerFunc.
func (h HandlerFunc) Serve(w ResponseWriter, r *Request) error {
	return h(w, r)
}

// HandleFunc is a convenience function to wrap a function into a HandlerFunc.
func HandleFunc(f func(ResponseWriter, *Request) error) HandlerFunc {
	return HandlerFunc(f)
}
