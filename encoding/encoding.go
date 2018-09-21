package encoding

import (
	"io"

	"github.com/ulule/ancre/timestamp"
)

type Encoder func(w io.Writer, t *timestamp.Timestamp, s *timestamp.Step) error

func Encode(w io.Writer, t *timestamp.Timestamp, fn Encoder) error {
	return encode(w, t, t.FirstStep, fn)
}

func encode(w io.Writer, t *timestamp.Timestamp, s *timestamp.Step, fn Encoder) error {
	err := fn(w, t, s)
	if err != nil {
		return err
	}

	for i := range s.Next {
		err = encode(w, t, s.Next[i], fn)
		if err != nil {
			return err
		}
	}
	return nil
}
