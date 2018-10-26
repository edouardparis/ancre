package ots

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/timestamp"
)

const RECURSION_LIMIT = 256
const URI_MAX_LENGTH = 1000

// MAX_OP_LENGTH is the maximum length of byte array
// used to describe operation.
const MAX_OP_LENGTH = 4096

func DecodeStep(ctx context.Context, r io.Reader, input []byte, currentTag *byte, attestations *[]attestation.Attestation) (*timestamp.Step, error) {
	if currentTag == nil {
		t, err := GetByte(r)
		if err != nil {
			return nil, err
		}
		currentTag = &t
	}

	next := []*timestamp.Step{}

	recursive := func(i []byte, t *byte) error {
		step, err := DecodeStep(ctx, r, i, t, attestations)
		if err != nil {
			if err != io.EOF {
				return err
			}

			return nil
		}

		next = append(next, step)

		return nil
	}

	switch *currentTag {
	case Attestation:
		a, err := DecodeAttestation(r, input)
		if err != nil {
			return nil, err
		}
		if attestations != nil {
			*attestations = append(*attestations, a)
		}
		return &timestamp.Step{Data: a}, nil

	case Fork:
		nextTag := byte(Fork)
		for nextTag == Fork {
			err := recursive(input, nil)
			if err != nil {
				return nil, err
			}

			nextTag, err = GetByte(r)
			if err != nil {
				return nil, err
			}
		}
		err := recursive(input, &nextTag)
		if err != nil {
			return nil, err
		}
		return &timestamp.Step{Output: input, Data: operation.NewFork(), Next: next}, nil

	default:
		op, err := DecodeOperation(r, *currentTag)
		if err != nil {
			return nil, err
		}
		output := op.Exec(input)
		err = recursive(output, nil)
		if err != nil {
			return nil, err
		}
		return &timestamp.Step{Data: op, Output: output, Next: next}, nil
	}
}

func DecodeOperationFromReader(r io.Reader) (*operation.Op, error) {
	t, err := GetByte(r)
	if err != nil {
		return nil, err
	}
	return DecodeOperation(r, t)
}

func DecodeOperation(r io.Reader, t byte) (*operation.Op, error) {
	switch t {
	case Sha1:
		return operation.NewOpSha1(), nil
	case Sha256:
		return operation.NewOpSha256(), nil
	case Keccak256:
		return operation.NewOpKeccak256(), nil
	case Ripemd160:
		return operation.NewOpRipemd160(), nil
	case Append:
		return DecodeOpAppend(r)
	case Prepend:
		return DecodeOpPrepend(r)
	}
	return nil, errors.New("operation does not exist")
}

// DecodeOpAppend append a byte array to the byte array.
func DecodeOpAppend(r io.Reader) (*operation.Op, error) {
	data, err := readData(r)
	if err != nil {
		return nil, err
	}

	return operation.NewOpAppend(data), nil
}

// DecodeOpPrepend prepend a byte array to the byte array.
func DecodeOpPrepend(r io.Reader) (*operation.Op, error) {
	data, err := readData(r)
	if err != nil {
		return nil, err
	}

	return operation.NewOpPrepend(data), nil
}

func DecodeAttestation(r io.Reader, input []byte) (attestation.Attestation, error) {
	t := make([]byte, ATTESTATION_SIZE_TAG)
	_, err := r.Read(t)
	if err != nil {
		return nil, err
	}

	len, err := ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(t, BITCOIN_TAG) {
		return DecodeBitcoin(r, input)
	} else if bytes.Equal(t, PENDING_TAG) {
		return DecodePending(r, input)
	}

	return DecodeUnknow(r, t, len, input)
}

func DecodeBitcoin(r io.Reader, input []byte) (attestation.Attestation, error) {
	height, err := ReadUInt64(r)
	if err != nil {
		return nil, err
	}
	return attestation.NewBitcoin(height, input), nil
}

func DecodePending(r io.Reader, input []byte) (attestation.Attestation, error) {
	length, err := ReadUInt64(r)
	if err != nil {
		return nil, err
	}

	if length > URI_MAX_LENGTH {
		return nil, errors.New("Attestation pending has a too loog uri length")
	}

	uri := make([]byte, length)
	n, err := r.Read(uri)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return attestation.NewPending(string(uri[:n]), input), nil
}

func DecodeUnknow(r io.Reader, tag []byte, l uint64, input []byte) (attestation.Attestation, error) {
	data := make([]byte, l)
	_, err := r.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return attestation.NewUnknow(data, tag, input), nil
}

func readData(r io.Reader) ([]byte, error) {
	n, err := ReadUInt64(r)
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

func Decode(ctx context.Context, r io.Reader, startDigest []byte) (*timestamp.Timestamp, error) {
	var attestations []attestation.Attestation
	firstStep, err := DecodeStep(ctx, r, startDigest, nil, &attestations)
	if err != nil {
		return nil, err
	}
	return timestamp.New(startDigest, firstStep, attestations...), nil
}
