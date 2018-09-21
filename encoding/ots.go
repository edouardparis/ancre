package encoding

import (
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
	"github.com/ulule/ancre/timestamp"
)

func ToOTS(w io.Writer, t *timestamp.Timestamp, s *timestamp.Step) error {
	if s == nil {
		return nil
	}

	err := OTSEncodeStep(w, s)
	if err != nil {
		return err
	}
	return nil
}

func OTSEncodeStep(w io.Writer, s *timestamp.Step) error {
	if s == nil {
		return nil
	}

	t := []byte{}
	if s.Match(operation.Fork) {
		t = append(t, tag.Fork)
	} else if s.Match(operation.Sha256) {
		t = append(t, tag.Sha256)
	} else if s.Match(operation.Sha256) {
		t = append(t, tag.Ripemd160)
	} else if s.Match(operation.Prepend) || s.Match(operation.Append) {
		return otsEncodeBinaryOperation(w, s)
	} else if s.Match(operation.Attestation) {
		t = s.Data.(attestation.Attestation).Encode()
	} else {
		return nil
	}
	_, err := w.Write(t)
	if err != nil {
		return err
	}

	return nil
}

func otsEncodeBinaryOperation(w io.Writer, s *timestamp.Step) error {
	t := []byte{}
	if s.Match(operation.Append) {
		t = append(t, tag.Append)
	} else {
		t = append(t, tag.Prepend)
	}

	_, err := w.Write(t)
	if err != nil {
		return err
	}

	data := s.Exec(nil)
	_, err = tag.WriteUint64(w, uint64(len(data)))
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}
