package p2p

import "github.com/toqns/toqns/foundation/address"

// ResponseWriter provides functionality to answer to a request.
type ResponseWriter interface {
	// Write writes the provided data as payload to the response.
	// Returns the bytes written or an error.
	Write([]byte) (int, error)

	// WriteStatus registers the status code for the response.
	WriteStatus(int)

	// WriteStatusWithExplanation registers the status code for the response
	// with the provided status text.
	WriteStatusWithExplanation(int, string)
}

// Response represents the response to a request.
type Response struct {
	// From is the responsing node.
	From address.Address

	// To is the node to receive the response.
	To address.Address

	// StatusCode is a status code for the response.
	StatusCode int

	// Status is a text with an explanation for StatusCode.
	Status string

	// Payload is data to be passed with the response.
	Payload []byte
}

// Write processes the received data for the response.
//
// Will set StatusCode to StatusOK unless a status has been
// provided via WriteStatus or WriteStatusWithExplanation.
func (r *Response) Write(data []byte) (int, error) {
	if r.StatusCode == 0 {
		r.WriteStatus(StatusOK)
	}
	r.Payload = data
	return len(r.Payload), nil
}

// WriteStatus registers the status code for the response.
func (r *Response) WriteStatus(code int) {
	r.StatusCode = code
	r.Status = StatusText(code)
}

// WriteStatusWithExplanation registers the status code for the response
// with the provided status text.
func (r *Response) WriteStatusWithExplanation(code int, e string) {
	r.StatusCode = code
	r.Status = e
}
