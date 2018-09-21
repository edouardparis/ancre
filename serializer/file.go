package serializer

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/ulule/ancre/decoding"
	"github.com/ulule/ancre/encoding"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
	"github.com/ulule/ancre/timestamp"
)

// Header magic bytes
// Designed to be give the user some information in a hexdump, while being
// identified as 'data' by the file utility.
var HEADER_MAGIC = []byte("\x00OpenTimestamps\x00\x00Proof\x00\xbf\x89\xe2\xe8\x84\xe8\x92\x94")
var VERSION = uint64(1)

// TimestampFile is a info file.
type TimestampFile struct {
	Digest     []byte
	DigestType *operation.Op
	Timestamp  *timestamp.Timestamp
}

// Infos is the TimestampFile func for displaying Info
func (t *TimestampFile) Infos(w io.Writer) error {
	digestInfo := []byte(fmt.Sprintf("%s\n",
		hex.EncodeToString(t.Digest),
	))
	_, err := w.Write(digestInfo)
	if err != nil {
		return err
	}

	return encoding.Encode(w, t.Timestamp, encoding.ToTEXT)
}

func NewTimeStampFile(digest []byte, op *operation.Op) *TimestampFile {
	t := &timestamp.Timestamp{}
	return &TimestampFile{Digest: digest, DigestType: op, Timestamp: t}
}

// TimestampFileFromReader creates TimestampFile from io reader.
func TimestampFileFromReader(r io.Reader) (*TimestampFile, error) {
	err := ReadMagic(r)
	if err != nil {
		return nil, err
	}

	err = ReadVersion(r)
	if err != nil {
		return nil, err
	}

	digestType, err := decoding.OTSDecodeOperationFromReader(r)
	if err != nil {
		return nil, err
	}

	digest := make([]byte, digestType.Length())
	_, err = r.Read(digest)
	if err != nil {
		return nil, err
	}

	t, err := decoding.FromOTS(context.Background(), r, digest)
	if err != nil {
		return nil, err
	}

	return &TimestampFile{Timestamp: t, DigestType: digestType, Digest: digest}, nil
}

func TimestampFileToWriter(t *TimestampFile, w io.Writer) error {
	err := WriteMagic(w)
	if err != nil {
		return err
	}

	err = WriteVersion(w)
	if err != nil {
		return err
	}

	_, err = w.Write(t.DigestType.Encode())
	if err != nil {
		return err
	}

	_, err = w.Write(t.Digest)
	if err != nil {
		return err
	}
	return encoding.Encode(w, t.Timestamp, encoding.ToOTS)
}

// ReadMagic reads the magic bytes and return a error
// if it is not the case.
func ReadMagic(r io.Reader) error {
	header := make([]byte, len(HEADER_MAGIC))
	_, err := r.Read(header)
	if err != nil {
		return err
	}

	if bytes.Equal(header, HEADER_MAGIC) {
		return nil
	}

	return errors.New("File has no magic header")
}

// WriteMagic writes the magic header
func WriteMagic(w io.Writer) error {
	_, err := w.Write(HEADER_MAGIC)
	return err
}

// ReadVersion reads the ots version and return a error
// if it is not the case.
func ReadVersion(r io.Reader) error {
	version, err := tag.ReadUInt64(r)
	if err != nil {
		return err
	}

	if version == VERSION {
		return nil
	}

	return errors.New("File has not the good version")
}

// WriteVersion writes the version
func WriteVersion(w io.Writer) error {
	_, err := w.Write([]byte{byte(VERSION)})
	return err
}
