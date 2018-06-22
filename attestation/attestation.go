package attestation

import (
	"bytes"
	"io"

	"github.com/ulule/ancre/tag"
)

const ATTESTATION_SIZE_TAG = 8
const URI_MAX_LENGTH = 1000

var BITCOIN_TAG = []byte("\x05\x88\x96\x0d\x73\xd7\x19\x01")
var PENDING_TAG = []byte("\x83\xdf\xe3\x0d\x2e\xf9\x0c\x8e")

type Attestation interface {
	Encode() []byte
	Match([]byte) bool
	Input() []byte
	Data() map[string]interface{}
}

func Decode(r io.Reader, input []byte) (Attestation, error) {
	t := make([]byte, ATTESTATION_SIZE_TAG)
	_, err := r.Read(t)
	if err != nil {
		return nil, err
	}

	len, err := tag.ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(t, BITCOIN_TAG) {
		return NewAttestBitcoin(r, input)
	} else if bytes.Equal(t, PENDING_TAG) {
		return NewAttestPending(r, input)
	}

	return NewAttestUnknow(r, t, len, input)
}
