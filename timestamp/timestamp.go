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
	Data() StepData
	Output() []byte
	HasNext() bool
}

type StepData interface {
	Match([]byte) bool
	Encode() []byte
}

type step struct {
	data   StepData
	output []byte
	next   []*step
}

func (s step) Data() StepData {
	return s.data
}

func (s step) Output() []byte {
	return s.output
}

func (s step) Encode() []byte {
	return s.data.Encode()
}

func (s step) Match(b []byte) bool {
	return s.data.Match(b)
}

func (s step) HasNext() bool {
	return len(s.next) > 0
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
		return &step{data: a}, nil

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
		return &step{output: input, data: operation.Fork{}, next: next}, nil

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
		return &step{data: op, output: output, next: next}, nil
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

type Encoder func(t *Timestamp, s Step) error

func Encode(t *Timestamp, fn Encoder) error {
	return encode(t, t.FirstStep, fn)
}

func encode(t *Timestamp, s *step, fn Encoder) error {
	err := fn(t, s)
	if err != nil {
		return err
	}

	for i := range s.next {
		err = encode(t, s.next[i], fn)
	}
	return nil
}
