package p2p

// RequestEncoder marshals data from the chosen encoding format.
type RequestEncoder interface {
	Marshal(any) ([]byte, error)
}

// RequestEncoderFunc wraps a RequestEncoder into a function.
type RequestEncoderFunc func(any) ([]byte, error)

// Marshal implements the RequestEncoder interface for RequestEncoderFunc.
func (r RequestEncoderFunc) Marshal(v any) ([]byte, error) {
	return r(v)
}

// RequestDecoder unmarshals data into the encoding format.
type RequestDecoder interface {
	Unmarshal([]byte, any) error
}

// RequestDecoderFunc wraps a RequestDecoder into a function.
type RequestDecoderFunc func([]byte, any) error

// Unmarshal implements the RequestDecoder interface for RequestDecoderFunc.
func (r RequestDecoderFunc) Unmarshal(data []byte, v any) error {
	return r(data, v)
}
