package attestation

import (
	"bytes"
	"io"

	"github.com/ulule/ancre/tag"
)

const ATTESTATION_SIZE_TAG = 8
const URI_MAX_LENGTH = 1000

type Attestation interface {
	Encode() []byte
	Match(int) bool
	Input() []byte
	Data() map[string]interface{}
	Exec([]byte) []byte
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

	if bytes.Equal(t, tag.BITCOIN_TAG) {
		return NewAttestBitcoin(r, input)
	} else if bytes.Equal(t, tag.PENDING_TAG) {
		return NewAttestPending(r, input)
	}

	return NewAttestUnknow(r, t, len, input)
}
