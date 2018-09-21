package operation

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ulule/ancre/tag"
)

func TestNewOpAppend(t *testing.T) {
	is := assert.New(t)

	b, err := hex.DecodeString(`f00688ac00000000`)
	is.NoError(err)

	r := bytes.NewReader(b)

	ta, err := tag.GetByte(r)
	is.NoError(err)

	is.Equal(ta, tag.Append)

	op, err := NewOpAppend(r)
	is.NoError(err)

	is.True(bytes.Equal(b, op.Tag))
}
