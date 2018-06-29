package timestamp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
)

func Information(w io.Writer, t *Timestamp) error {
	return DisplayStep(w, t, &t.FirstStep)
}

func DisplayStep(w io.Writer, t *Timestamp, step *Step) error {
	if step.Next == nil {
		attestation, ok := step.Data.(attestation.Attestation)
		if !ok {
			return errors.New("step is not an attestion")
		}
		return DisplayAttestation(w, attestation)
	} else {
		err := DisplayOperation(w, step)
		if err != nil {
			return err
		}
	}

	for i := range step.Next {
		err := DisplayStep(w, t, &step.Next[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func DisplayOperation(w io.Writer, step *Step) error {
	op, ok := step.Data.(operation.Operation)
	if !ok {
		return errors.New("step is not an operation")
	}

	_, err := w.Write([]byte("execute "))
	if err != nil {
		return err
	}

	// Temporary return nil
	if step.Data == nil {
		return nil
	}

	operations := map[byte]string{
		tag.Sha256:    "SHA256()\n",
		tag.Ripemd160: "RIPEMD160()\n",
	}

	for t := range operations {
		if op.Match([]byte{t}) {
			_, err = w.Write([]byte(operations[t]))
			break
		}
	}

	if op.Match([]byte{tag.Append}) {
		argument := hex.EncodeToString(op.Exec(nil))
		_, err = w.Write([]byte(fmt.Sprintf("Append(%s)\n", argument)))
	}

	if op.Match([]byte{tag.Prepend}) {
		argument := hex.EncodeToString(op.Exec(nil))
		_, err = w.Write([]byte(fmt.Sprintf("Prepend(%s)\n", argument)))
	}

	if err != nil {
		return err
	}

	result := fmt.Sprintf("result: %s\n", hex.EncodeToString(step.Output))
	_, err = w.Write([]byte(result))
	return err
}

func DisplayAttestation(w io.Writer, attest attestation.Attestation) error {
	_, err := w.Write([]byte("attested by :\n"))
	if err != nil {
		return err
	}

	for k, v := range attest.Data() {
		_, err := w.Write([]byte(fmt.Sprintf("# %s : %v\n", k, v)))
		if err != nil {
			return err
		}
	}

	_, err = w.Write([]byte("\n"))
	return err
}
