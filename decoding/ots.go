package decoding

import (
	"context"
	"errors"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
	"github.com/ulule/ancre/timestamp"
)

const RECURSION_LIMIT = 256

// MAX_OP_LENGTH is the maximum length of byte array
// used to describe operation.
const MAX_OP_LENGTH = 4096

func OTSDecodeStep(ctx context.Context, r io.Reader, input []byte, currentTag *byte, attestations []attestation.Attestation) (*timestamp.Step, error) {
	if currentTag == nil {
		t, err := tag.GetByte(r)
		if err != nil {
			return nil, err
		}
		currentTag = &t
	}

	next := []*timestamp.Step{}

	recursive := func(i []byte, t *byte) error {
		step, err := OTSDecodeStep(ctx, r, i, t, attestations)
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
	case tag.Attestation:
		a, err := attestation.Decode(r, input)
		if err != nil {
			return nil, err
		}
		attestations = append(attestations, a)
		return &timestamp.Step{Data: a}, nil

	case tag.Fork:
		nextTag := byte(tag.Fork)
		for nextTag == tag.Fork {
			err := recursive(input, nil)
			if err != nil {
				return nil, err
			}

			nextTag, err = tag.GetByte(r)
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
		op, err := OTSDecodeOperation(r, *currentTag)
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

func OTSDecodeOperationFromReader(r io.Reader) (*operation.Op, error) {
	t, err := tag.GetByte(r)
	if err != nil {
		return nil, err
	}
	return OTSDecodeOperation(r, t)
}

func OTSDecodeOperation(r io.Reader, t byte) (*operation.Op, error) {
	switch t {
	case tag.Sha256:
		return operation.NewOpSha256(), nil
	case tag.Ripemd160:
		return operation.NewOpRipemd160(), nil
	case tag.Append:
		return DecodeOpAppend(r)
	case tag.Prepend:
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

func FromOTS(ctx context.Context, r io.Reader, startDigest []byte) (*timestamp.Timestamp, error) {
	var attestations []attestation.Attestation
	firstStep, err := OTSDecodeStep(ctx, r, startDigest, nil, attestations)
	if err != nil {
		return nil, err
	}
	return timestamp.New(firstStep, attestations...), nil
}
