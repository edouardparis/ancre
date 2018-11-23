package ots

import (
	"bytes"
	"io"

	"github.com/ulule/ancre/attestation"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/timestamp"
)

func Encode(w io.Writer, t *timestamp.Timestamp) error {
	return encode(w, t, t.FirstStep)
}

func encode(w io.Writer, t *timestamp.Timestamp, s *timestamp.Step) error {
	if s == nil {
		return nil
	}

	if !s.Match(operation.Fork) {
		err := EncodeStep(w, s)
		if err != nil {
			return err
		}
	}

	for i := range s.Next {
		if s.Match(operation.Fork) && i+1 != len(s.Next) {
			_, err := w.Write([]byte{Fork})
			if err != nil {
				return err
			}
		}

		err := encode(w, t, s.Next[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func EncodeStep(w io.Writer, s *timestamp.Step) error {
	if s == nil {
		return nil
	}

	t := []byte{}
	if s.Match(operation.Sha256) {
		t = append(t, Sha256)
	} else if s.Match(operation.Ripemd160) {
		t = append(t, Ripemd160)
	} else if s.Match(operation.Prepend) || s.Match(operation.Append) {
		return EncodeBinaryOperation(w, s)
	} else if s.Match(operation.Attestation) {
		t = EncodeAttestation(s.Data.(attestation.Attestation))
	} else {
		return nil
	}
	_, err := w.Write(t)
	if err != nil {
		return err
	}

	return nil
}

func EncodeBinaryOperation(w io.Writer, s *timestamp.Step) error {
	t := []byte{}
	if s.Match(operation.Append) {
		t = append(t, Append)
	} else {
		t = append(t, Prepend)
	}

	_, err := w.Write(t)
	if err != nil {
		return err
	}

	data := s.Exec(nil)
	_, err = WriteUint64(w, uint64(len(data)))
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func EncodeAttestation(a attestation.Attestation) []byte {
	if b, ok := a.(*attestation.Bitcoin); ok {
		return EncodeBitcoin(b)
	}

	if p, ok := a.(*attestation.Pending); ok {
		return EncodePending(p)
	}

	if u, ok := a.(*attestation.Unknow); ok {
		return EncodeUnknow(u)
	}

	return nil
}

func EncodeBitcoin(b *attestation.Bitcoin) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{Attestation})
	buf.Write(BITCOIN_TAG)

	tmp := new(bytes.Buffer)
	n, _ := WriteUint64(tmp, b.Height())

	WriteUint64(buf, uint64(n))
	buf.Write(tmp.Bytes())
	return buf.Bytes()
}

// EncodePending encodes the Attestation
func EncodePending(p *attestation.Pending) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte{Attestation})
	buf.Write(PENDING_TAG)

	b := []byte(p.Uri())

	WriteUint64(buf, uint64(len(b)))
	buf.Write(b)
	return buf.Bytes()
}

func EncodeUnknow(u *attestation.Unknow) []byte {
	return append(u.Tag(), u.RawData()...)
}
