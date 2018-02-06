package json

import (
	"encoding/json"
	"io"

	"github.com/hermes/ghord/codec"
)

type Codec struct{}

func NewCodec() Codec {
	return Codec{}
}

func (c Codec) Name() string {
	return "json"
}

func (c Codec) NewEncoder(w io.Writer) codec.Encoder {
	return json.NewEncoder(w)
}

func (c Codec) NewDecoder(r io.Reader) codec.Decoder {
	return json.NewDecoder(r)
}
