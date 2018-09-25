package attestation

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
)

type Bitcoin struct {
	height uint64
	input  []byte
}

func (b Bitcoin) Encode() []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{tag.Attestation})
	buf.Write(tag.BITCOIN_TAG)

	tmp := new(bytes.Buffer)
	n, _ := tag.WriteUint64(tmp, b.height)

	tag.WriteUint64(buf, uint64(n))
	buf.Write(tmp.Bytes())
	return buf.Bytes()
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

func (b Bitcoin) Data() map[string]interface{} {
	return map[string]interface{}{
		"input":  hex.EncodeToString(b.input),
		"height": b.height,
		"entity": "Bitcoin (BTC)",
	}
}

func NewAttestBitcoin(r io.Reader, input []byte) (Attestation, error) {
	height, err := tag.ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	return &Bitcoin{height: height, input: input}, nil
}
