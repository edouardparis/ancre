package timestamp

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
)

const RECURSION_LIMIT = 256

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

// Timstamp is the main structure representing a timestamp.
type Timestamp struct {
	Digest       []byte
	FirstStep    *Step
	Attestations []attestation.Attestation
}

// Sub creates a timestamp with the given digest and a subset
// of the parent timestamp.
func (t *Timestamp) Sub(digest []byte) (*Timestamp, error) {
	return nil, nil
}

func (t *Timestamp) Merge(ts *Timestamp) error {
	if ts == nil || ts.FirstStep == nil {
		return nil
	}

	t.Attestations = append(t.Attestations, ts.Attestations...)

	if t.FirstStep == nil {
		t.FirstStep = ts.FirstStep
		return nil
	}

	if bytes.Equal(t.Digest, ts.Digest) {
		if t.FirstStep.Match(operation.Fork) {
			t.FirstStep.Next = append(t.FirstStep.Next, ts.FirstStep)
		} else if ts.FirstStep != nil {
			firstStep := *t.FirstStep
			t.FirstStep = &Step{
				Output: t.Digest,
				Data:   operation.NewFork(),
				Next:   []*Step{&firstStep, ts.FirstStep},
			}
		}
	} else {
		ok := mergeStep(t.FirstStep, ts.Digest, ts.FirstStep)
		if !ok {
			return fmt.Errorf(
				"Failed to merge timestamp with digest: %s in timestamp with digest %s",
				hex.EncodeToString(ts.Digest), hex.EncodeToString(t.Digest),
			)
		}
	}

	return nil
}

// TODO: some concurrency stuff would be nice
func mergeStep(s *Step, digest []byte, first *Step) bool {
	if bytes.Equal(s.Output, digest) {
		s.Next = append(s.Next, first)
		return true
	}

	for i := range s.Next {
		ok := mergeStep(s.Next[i], digest, first)
		if ok {
			return ok
		}
	}
	return false
}

// New returns a pointer to a new timestamp instance.
func New(digest []byte, firstStep *Step, attestations ...attestation.Attestation) *Timestamp {
	return &Timestamp{Digest: digest, FirstStep: firstStep, Attestations: attestations}
}
