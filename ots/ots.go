package ots

import (
	"errors"
	"io"
)

const (
	// Operations
	Sha256    = byte(0x08)
	Ripemd160 = byte(0x03)
	Append    = byte(0xf0)
	Prepend   = byte(0xf1)

	// Step
	Attestation = byte(0x00)
	Fork        = byte(0xff)

	ATTESTATION_SIZE_TAG = 8
)

var BITCOIN_TAG = []byte("\x05\x88\x96\x0d\x73\xd7\x19\x01")
var PENDING_TAG = []byte("\x83\xdf\xe3\x0d\x2e\xf9\x0c\x8e")

// GetByte return one byte from io.Reader
func GetByte(r io.Reader) (byte, error) {
	t := make([]byte, 1)
	_, err := r.Read(t)
	if err != nil {
		return 0, err
	}

	return t[0], nil
}

func ReadUInt64(r io.Reader) (uint64, error) {
	var (
		x uint64
		s uint
	)

	for i := 0; i <= 9; {
		b, err := GetByte(r)
		if err != nil {
			return 0, err
		}

		if b < 0x80 {
			return x | uint64(b)<<s, nil
		}

		x |= uint64(b&0x7f) << s
		s += 7
	}

	return 0, errors.New("UInt64 overflow")
}

func WriteUint64(w io.Writer, i uint64) (int, error) {
	if i == 0 {
		return w.Write([]byte{0x00})
	}

	var total int
	for i > 0 {
		b := byte(i & 0x7f)
		if i > uint64(0x7f) {
			b |= 0x80
		}
		n, err := w.Write([]byte{b})
		if err != nil {
			return total, err
		}
		total += n
		if i <= 0x7f {
			break
		}
		i >>= 7
	}
	return total, nil
}
