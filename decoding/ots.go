package decoding

import (
	"context"
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
		return &timestamp.Step{Output: input, Data: operation.Fork{}, Next: next}, nil

	default:
		op, err := operation.Decode(r, *currentTag)
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

func FromOTS(ctx context.Context, r io.Reader, startDigest []byte) (*timestamp.Timestamp, error) {
	var attestations []attestation.Attestation
	firstStep, err := OTSDecodeStep(ctx, r, startDigest, nil, attestations)
	if err != nil {
		return nil, err
	}
	return timestamp.New(firstStep, attestations...), nil
}
