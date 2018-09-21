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
	Match(int) bool
	Exec([]byte) []byte
}

type Step struct {
	Data   StepData
	Output []byte
	Next   []*Step
}

func (s Step) Exec(input []byte) []byte {
	return s.Data.Exec(input)
}

func (s Step) Match(i int) bool {
	return s.Data.Match(i)
}

func (s Step) HasNext() bool {
	return len(s.Next) > 0
}

// New returns a pointer to a new timestamp instance.
func New(firstStep *Step, attestations ...attestation.Attestation) *Timestamp {
	return &Timestamp{FirstStep: firstStep, Attestations: attestations}
}
