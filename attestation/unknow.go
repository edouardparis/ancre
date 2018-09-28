package attestation

import (
	"encoding/hex"

	"github.com/ulule/ancre/operation"
)

type Unknow struct {
	tag   []byte
	data  []byte
	input []byte
}

func (u Unknow) Tag() []byte {
	return u.tag
}

func (u Unknow) RawData() []byte {
	return u.data
}

func (u Unknow) Encode() []byte {
	return append(u.tag, u.data...)
}

func (u Unknow) Match(i int) bool {
	return i == operation.Attestation
}

func (u Unknow) Input() []byte {
	return u.input
}

func (u Unknow) Exec(input []byte) []byte {
	return input
}

func (u Unknow) Data() map[string]interface{} {
	return map[string]interface{}{
		"input":  hex.EncodeToString(u.input),
		"data":   hex.EncodeToString(u.data),
		"entity": "unknow",
	}
}

func NewUnknow(data []byte, tag []byte, input []byte) *Unknow {
	return &Unknow{data: data, tag: tag, input: input}
}
