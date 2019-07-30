package attestation

import (
	"encoding/hex"

	"github.com/edouardparis/ancre/operation"
)

type Bitcoin struct {
	height uint64
	input  []byte
}

func (b Bitcoin) Match(i int) bool {
	return i == operation.Attestation
}

func (b Bitcoin) Exec(input []byte) []byte {
	return input
}

func (b Bitcoin) Input() []byte {
	return b.input
}

func (b Bitcoin) Height() uint64 {
	return b.height
}

func (b Bitcoin) Data() map[string]interface{} {
	return map[string]interface{}{
		"input":  hex.EncodeToString(b.input),
		"height": b.Height,
		"entity": "Bitcoin (BTC)",
	}
}

func NewBitcoin(h uint64, input []byte) *Bitcoin {
	return &Bitcoin{height: h, input: input}
}
