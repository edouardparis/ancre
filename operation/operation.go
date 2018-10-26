package operation

import (
	"crypto/sha1"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const (
	Fork = iota + 1
	Sha1
	Sha256
	Keccak256
	Ripemd160
	Append
	Prepend
	Attestation
)

type Operation interface {
	Match(int) bool
	Exec(input []byte) []byte
}

// Op is the type of an operation done on a timestamp.
type Op struct {
	Kind         int
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

func (o Op) Match(i int) bool {
	return o.Kind == i
}

// NewOpSha1 is the sha1 operation checkSum.
func NewOpSha1() *Op {
	return &Op{Sha1, 20, func(data []byte) []byte {
		b := sha1.Sum(data)
		return b[:]
	}}
}

// NewOpSha256 is the sha256 operation checkSum.
func NewOpSha256() *Op {
	return &Op{Sha256, 32, func(data []byte) []byte {
		b := sha256.Sum256(data)
		return b[:]
	}}
}

// NewOpKeccak256 is the sha3.Sum256 operation checkSum.
func NewOpKeccak256() *Op {
	return &Op{Keccak256, 32, func(data []byte) []byte {
		b := sha3.Sum256(data)
		return b[:]
	}}
}

// NewOpRipemd160 is the sha256 operation checkSum.
// golang/x/crypto Ripemd160 does not return error,
// Output of hash.Write can be ignored.
func NewOpRipemd160() *Op {
	return &Op{Ripemd160, 20, func(data []byte) []byte {
		hash := ripemd160.New()
		hash.Write(data)
		return hash.Sum(nil)
	}}
}

// NewOpAppend append a byte array to the byte array.
func NewOpAppend(data []byte) *Op {
	return &Op{Append, 32, func(input []byte) []byte {
		return append(input, data...)
	}}
}

// NewOpPrepend prepend a byte array to the byte array.
func NewOpPrepend(data []byte) *Op {
	return &Op{Prepend, 32, func(input []byte) []byte {
		return append(data, input...)
	}}
}

func NewFork() *Op {
	return &Op{Fork, 0, func(input []byte) []byte {
		return input
	}}
}
