package format

import (
	"io"

	"github.com/ulule/ancre/timestamp"
)

func ToOTS(w io.Writer) timestamp.Encoder {
	return func(t *timestamp.Timestamp, s timestamp.Step) error {
		if s == nil {
			return nil
		}

		_, err := w.Write(s.Encode())
		if err != nil {
			return err
		}
		return nil
	}
}
