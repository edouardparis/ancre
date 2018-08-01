package format

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

func ToTEXT(w io.Writer) func(timestamp.Step) error {
	return func(s timestamp.Step) error {
		if s == nil {
			return nil
		}

		if s.Match(attestation.BITCOIN_TAG) || s.Match(attestation.PENDING_TAG) {
			attest, ok := s.GetData().(attestation.Attestation)
			if !ok {
				return errors.New("step is not an attestion")
			}
			return DisplayAttestation(w, attest)
		} else if s.Match([]byte{tag.Fork}) {
			_, err := w.Write([]byte("(fork two ways)\n"))
			if err != nil {
				return err
			}
		} else {
			op, ok := s.GetData().(operation.Operation)
			if !ok {
				return errors.New("step is not an operation")
			}
			err := DisplayOperation(w, op, s.GetOutput())
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func DisplayOperation(w io.Writer, op operation.Operation, output []byte) error {
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
