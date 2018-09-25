package attestation

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"

	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
)

type Pending struct {
	uri   string
	input []byte
}

// Encode encodes the Attestation
func (p Pending) Encode() []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{tag.Attestation})
	buf.Write(tag.PENDING_TAG)

	tmp := new(bytes.Buffer)
	n, _ := tmp.WriteString(p.uri)

	tag.WriteUint64(buf, uint64(n))
	buf.Write(tmp.Bytes())
	return buf.Bytes()
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

func NewAttestPending(r io.Reader, input []byte) (Attestation, error) {
	length, err := tag.ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	if length > URI_MAX_LENGTH {
		return nil, errors.New("Attestation pending has a too loog uri length")
	}

	uri := make([]byte, length)
	_, err = r.Read(uri)
	if err != nil {
		return nil, err
	}

	return &Pending{uri: string(uri), input: input}, nil
}
