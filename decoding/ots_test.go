package decoding

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/tag"
)

func TestDecodeOpAppend(t *testing.T) {
	is := assert.New(t)

	b, err := hex.DecodeString(`f00688ac00000000`)
	is.NoError(err)

	r := bytes.NewReader(b)

	ta, err := tag.GetByte(r)
	is.NoError(err)

	is.Equal(ta, tag.Append)

	op, err := DecodeOpAppend(r)
	is.NoError(err)

	is.True(op.Match(operation.Append))
}
