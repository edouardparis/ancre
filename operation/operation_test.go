package operation

import (
	"bytes"
	"encoding/hex"
	"testing"
)

const hexDump = `f00688ac00000000`

func TestNewOpAppend(t *testing.T) {
	b, err := hex.DecodeString(`f00688ac00000000`)
	if err != nil {
		t.Errorf("Expected no error: \n %s", err.Error())
	}

	r := bytes.NewReader(b)
	op, err := NewOpAppend(r)
	if err != nil {
		t.Errorf("Expected no error: \n %s", err.Error())
	}

	if !bytes.Equal(b, op.Tag) {
		t.Errorf("Expected equal tags")
	}
}
