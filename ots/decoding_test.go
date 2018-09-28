package ots

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ulule/ancre/operation"
)

func TestDecodeOpAppend(t *testing.T) {
	is := assert.New(t)

	b, err := hex.DecodeString(`f00688ac00000000`)
	is.NoError(err)

	r := bytes.NewReader(b)

	ta, err := GetByte(r)
	is.NoError(err)

	is.Equal(ta, Append)

	op, err := DecodeOpAppend(r)
	is.NoError(err)

	is.True(op.Match(operation.Append))
}
