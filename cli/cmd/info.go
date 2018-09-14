package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/ulule/ancre/serializer"
)

// Info open the given file and display information.
func Info(filepath string) error {
	if path.Ext(filepath) != ".ots" {
		return fmt.Errorf("%s is not an ots file", filepath)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	ots, err := serializer.TimestampFileFromReader(file)
	if err != nil {
		return err
	}

	return ots.Infos(os.Stderr)
}
