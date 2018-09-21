package operation

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/ripemd160"

	"github.com/ulule/ancre/tag"
)

// MAX_OP_LENGTH is the maximum length of byte array
// used to describe operation.
const MAX_OP_LENGTH = 4096

const (
	Fork = iota + 1
	Sha256
)

type Operation interface {
	Encode() []byte
	Match([]byte) bool
	Exec(input []byte) []byte
}

// Op is the type of an operation done on a timestamp.
type Op struct {
	Tag          []byte
	DigestLength int
	Apply        func([]byte) []byte
}

// Exec execute the operation
func (o Op) Exec(input []byte) []byte {
	return o.Apply(input)
}

// Length returns the operation lenght
func (o Op) Length() int {
	return o.DigestLength
}

// Encode serializes the Op
func (o Op) Encode() []byte {
	return o.Tag
}

func (o Op) Match(b []byte) bool {
	return len(o.Tag) != 0 && bytes.Equal(o.Tag[:1], b)
}

// NewOpSha256 is the sha256 operation checkSum.
func NewOpSha256() *Op {
	return &Op{[]byte{tag.Sha256}, 32, func(data []byte) []byte {
		b := sha256.Sum256(data)
		return b[:]
	}}
}

// NewOpRipemd160 is the sha256 operation checkSum.
// golang/x/crypto Ripemd160 does not return error,
// Output of hash.Write can be ignored.
func NewOpRipemd160() *Op {
	return &Op{[]byte{tag.Ripemd160}, 20, func(data []byte) []byte {
		hash := ripemd160.New()
		hash.Write(data)
		return hash.Sum(nil)
	}}
}

// NewOpAppend append a byte array to the byte array.
func NewOpAppend(r io.Reader) (*Op, error) {
	data, err := readData(r)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	_, err = b.Write([]byte{tag.Append})
	if err != nil {
		return nil, err
	}

	_, err = tag.WriteUint64(&b, uint64(len(data)))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(data)
	if err != nil {
		return nil, err
	}

	return &Op{b.Bytes(), 32, func(input []byte) []byte {
		return append(input, data...)
	}}, nil
}

// NewOpPrepend append a byte array at the end of the byte array.
func NewOpPrepend(r io.Reader) (*Op, error) {
	data, err := readData(r)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	_, err = b.Write([]byte{tag.Prepend})
	if err != nil {
		return nil, err
	}

	_, err = tag.WriteUint64(&b, uint64(len(data)))
	if err != nil {
		return nil, err
	}

	_, err = b.Write(data)
	if err != nil {
		return nil, err
	}

	return &Op{b.Bytes(), 32, func(input []byte) []byte {
		return append(data, input...)
	}}, nil
}

func readData(r io.Reader) ([]byte, error) {
	n, err := tag.ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	if n > MAX_OP_LENGTH {
		return nil, errors.New("Input too long for binary operation ")
	}

	data := make([]byte, n)
	_, err = r.Read(data)
	return data, err
}

func NewFork() *Op {
	return &Op{[]byte{tag.Fork}, 0, func(input []byte) []byte {
		return input
	}}
}
