package operation

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ulule/ancre/tag"
)

func TestNewOpAppend(t *testing.T) {
	b, err := hex.DecodeString(`f00688ac00000000`)
	if err != nil {
		t.Errorf("Expected no error: \n %s", err.Error())
	}

	r := bytes.NewReader(b)

	ta, err := tag.GetByte(r)
	if err != nil {
		t.Errorf("Expected no error: \n %s", err.Error())
	}

	if ta != tag.Append {
		t.Errorf("Expected equal tags")
	}

	op, err := NewOpAppend(r)
	if err != nil {
		t.Errorf("Expected no error: \n %s", err.Error())
	}

	if !bytes.Equal(b, op.Tag) {
		t.Errorf("Expected equal tags")
	}
}
