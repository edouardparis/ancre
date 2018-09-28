package attestation

import (
	"encoding/hex"

	"github.com/ulule/ancre/operation"
)

type Pending struct {
	uri   string
	input []byte
}

func (p Pending) Uri() string {
	return p.uri
}

// Deser deserializes the Attestation
func (p Pending) Match(i int) bool {
	return i == operation.Attestation
}

func (p Pending) Exec(input []byte) []byte {
	return input
}

func (p Pending) Input() []byte {
	return p.input
}

func (p Pending) Data() map[string]interface{} {
	return map[string]interface{}{
		"uri":    p.uri,
		"input":  hex.EncodeToString(p.input),
		"entity": "pending",
	}
}

func NewPending(uri string, input []byte) *Pending {
	return &Pending{uri: uri, input: input}
}
