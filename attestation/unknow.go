package attestation

import (
	"bytes"
	"encoding/hex"
	"io"
)

type Unknow struct {
	tag   []byte
	data  []byte
	input []byte
}

func (u Unknow) Encode() []byte {
	return append(u.tag, u.data...)
}

func (u Unknow) Match(b []byte) bool {
	return bytes.Equal(u.tag, b)
}

func (u Unknow) Input() []byte {
	return u.input
}

func (u Unknow) Data() map[string]interface{} {
	return map[string]interface{}{
		"input":  hex.EncodeToString(u.input),
		"data":   hex.EncodeToString(u.data),
		"entity": "unknow",
	}
}

func NewAttestUnknow(r io.Reader, tag []byte, len uint64, input []byte) (Attestation, error) {
	data := make([]byte, len)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}

	return &Unknow{data: data, tag: tag, input: input}, nil
}
