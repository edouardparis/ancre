package attestation

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"

	"github.com/ulule/ancre/tag"
)

type Bitcoin struct {
	height uint64
	input  []byte
}

func (b Bitcoin) Encode() []byte {
	buf := new(bytes.Buffer)
	buf.Write(BITCOIN_TAG)
	binary.Write(buf, binary.LittleEndian, b.height)
	return buf.Bytes()
}

func (b Bitcoin) Match(a []byte) bool {
	return bytes.Equal(BITCOIN_TAG, a)
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
