package codec

import "io"

// Standard encoder interface, usually created from an io.Writer
// Important: Needs to be stateless, in that each call to Encode() must
// act the same as any other call
type Encoder interface {
	// Encode given interface, or error
	Encode(v interface{}) error
}

// Standatd decoder interface, usually created from an io.Reader
// Important: Needs to be stateless, in that each call to Decode() must
// act the same as any other call
type Decoder interface {
	// Decode into the given interface, or error
	Decode(v interface{}) error
}

// A codec for coding ghord messages
type Codec interface {
	Name() string
	NewEncoder(io.Writer) Encoder
	NewDecoder(io.Reader) Decoder
}
