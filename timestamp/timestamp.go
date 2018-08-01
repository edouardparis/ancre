package timestamp

import (
	"context"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
)

const RECURSION_LIMIT = 256

type Timestamp struct {
	FirstStep    Step
	Attestations []attestation.Attestation
}

type StepData interface {
	Encode() []byte
	Match([]byte) bool
}

type Step struct {
	Data   StepData
	Output []byte
	Next   []Step
}

func (ts *Timestamp) DecodeStep(ctx context.Context, r io.Reader, input []byte, currentTag *byte) (Step, error) {
	if currentTag == nil {
		t, err := tag.GetByte(r)
		if err != nil {
			return Step{}, err
		}
		currentTag = &t
	}

	next := []Step{}

	recursive := func(i []byte, t *byte) error {
		step, err := ts.DecodeStep(ctx, r, i, t)
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
			return Step{}, err
		}
		ts.Attestations = append(ts.Attestations, a)
		return Step{Data: a}, nil

	case tag.Fork:
		nextTag := byte(tag.Fork)
		for nextTag == tag.Fork {
			err := recursive(input, nil)
			if err != nil {
				return Step{}, err
			}

			nextTag, err = tag.GetByte(r)
			if err != nil {
				return Step{}, err
			}
		}
		err := recursive(input, &nextTag)
		if err != nil {
			return Step{}, err
		}
		return Step{Output: input, Next: next}, nil

	default:
		op, err := operation.Decode(r, *currentTag)
		if err != nil {
			return Step{}, err
		}
		output := op.Exec(input)
		err = recursive(output, nil)
		if err != nil {
			return Step{}, err
		}
		return Step{Data: op, Output: output, Next: next}, nil
	}
}

func (t *Timestamp) Decode(ctx context.Context, r io.Reader, startDigest []byte) error {
	firstStep, err := t.DecodeStep(ctx, r, startDigest, nil)
	if err != nil {
		return err
	}

	t.FirstStep = firstStep

	return nil
}

func EncodeRecursive(s Step, w io.Writer) error {
	_, err := w.Write(s.Data.Encode())
	if err != nil {
		return err
	}

	if len(s.Next) > 1 {
		for i := 0; i < len(s.Next); i++ {
			_, err := w.Write([]byte{tag.Fork})
			if err != nil {
				return err
			}
		}
	}

	for i := range s.Next {
		err := EncodeRecursive(s.Next[i], w)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Timestamp) Encode(w io.Writer) error {
	return EncodeRecursive(t.FirstStep, w)
}
