package timestamp

import (
	"github.com/ulule/ancre/attestation"
)

const RECURSION_LIMIT = 256

type Timestamp struct {
	FirstStep    *Step
	Attestations []attestation.Attestation
}

type StepData interface {
	Match([]byte) bool
	Encode() []byte
}

type Step struct {
	Data   StepData
	Output []byte
	Next   []*Step
}

func (s Step) Encode() []byte {
	return s.Data.Encode()
}

func (s Step) Match(b []byte) bool {
	return s.Data.Match(b)
}

func (s Step) HasNext() bool {
	return len(s.Next) > 0
}

type Encoder func(t *Timestamp, s *Step) error

func Encode(t *Timestamp, fn Encoder) error {
	return encode(t, t.FirstStep, fn)
}

func encode(t *Timestamp, s *Step, fn Encoder) error {
	err := fn(t, s)
	if err != nil {
		return err
	}

	for i := range s.Next {
		err = encode(t, s.Next[i], fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// New returns a pointer to a new timestamp instance.
func New(firstStep *Step, attestations ...attestation.Attestation) *Timestamp {
	return &Timestamp{FirstStep: firstStep, Attestations: attestations}
}
