package encoding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
	"github.com/ulule/ancre/timestamp"
)

func ToTEXT(w io.Writer, t *timestamp.Timestamp, s *timestamp.Step) error {
	if s == nil {
		return nil
	}

	if !s.HasNext() {
		attest, ok := s.Data.(attestation.Attestation)
		if !ok {
			return errors.New("step is not an attestion")
		}
		return attestationToText(w, attest)
	} else if s.Match([]byte{tag.Fork}) {
		_, err := w.Write([]byte("(fork two ways)\n"))
		if err != nil {
			return err
		}
	} else {
		op, ok := s.Data.(operation.Operation)
		if !ok {
			return errors.New("step is not an operation")
		}
		err := operationToText(w, op, s.Output)
		if err != nil {
			return err
		}
	}
	return nil
}

func operationToText(w io.Writer, op operation.Operation, output []byte) error {
	_, err := w.Write([]byte("execute "))
	if err != nil {
		return err
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

	result := fmt.Sprintf("result: %s\n", hex.EncodeToString(output))
	_, err = w.Write([]byte(result))
	return err
}

func attestationToText(w io.Writer, attest attestation.Attestation) error {
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