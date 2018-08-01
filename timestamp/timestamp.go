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
	FirstStep    *step
	Attestations []attestation.Attestation
}

type Step interface {
	StepData
	GetData() StepData
	GetOutput() []byte
}

type StepData interface {
	Match([]byte) bool
	Encode() []byte
}

type step struct {
	Data   StepData
	Output []byte
	Next   []*step
}

func (s step) GetData() StepData {
	return s.Data
}

func (s step) GetOutput() []byte {
	return s.Output
}

func (s step) Encode() []byte {
	return s.Data.Encode()
}

func (s step) Match(b []byte) bool {
	return s.Data.Match(b)
}

func (ts *Timestamp) DecodeStep(ctx context.Context, r io.Reader, input []byte, currentTag *byte) (*step, error) {
	if currentTag == nil {
		t, err := tag.GetByte(r)
		if err != nil {
			return nil, err
		}
		currentTag = &t
	}

	next := []*step{}

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
			return nil, err
		}
		ts.Attestations = append(ts.Attestations, a)
		return &step{Data: a}, nil

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
		return &step{Output: input, Data: operation.Fork{}, Next: next}, nil

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
		return &step{Data: op, Output: output, Next: next}, nil
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

func Encode(t *Timestamp, fn func(s Step) error) error {
	return encode(t.FirstStep, fn)
}

func encode(s *step, fn func(s Step) error) error {
	err := fn(s)
	if err != nil {
		return err
	}

	for i := range s.Next {
		err = encode(s.Next[i], fn)
	}
	return nil
}
