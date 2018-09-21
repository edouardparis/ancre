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
		return operation.NewOpAppend(r)
	case tag.Prepend:
		return operation.NewOpPrepend(r)
	}
	return nil, errors.New("operation does not exist")
}

func FromOTS(ctx context.Context, r io.Reader, startDigest []byte) (*timestamp.Timestamp, error) {
	var attestations []attestation.Attestation
	firstStep, err := OTSDecodeStep(ctx, r, startDigest, nil, attestations)
	if err != nil {
		return nil, err
	}
	return timestamp.New(firstStep, attestations...), nil
}
