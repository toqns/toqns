package p2p

import "github.com/toqns/toqns/foundation/address"

// Request represents a request to a node.
type Request struct {
	// To is the address the request is intended for.
	To address.Address

	// From is the address of the sending host.
	From address.Address

	// Payload is a slice of bytes representing the request's payload.
	Payload []byte

	// Response holds the response to the request.
	Response *Response
}
