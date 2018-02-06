package gob

import (
	"encoding/gob"
	"io"
)

type Codec struct{}

func NewCodec() Codec {
	return Codec{}
}

func (c Codec) NewEncoder(w io.Writer) *gob.Encoder {
	return gob.NewEncoder(w)
}

func (c Codec) NewDecoder(r io.Reader) *gob.Decoder {
	return gob.NewDecoder(r)
}
