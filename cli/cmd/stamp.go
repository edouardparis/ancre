package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ulule/ancre/calendar/client"
	"github.com/ulule/ancre/logging"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/serializer"
)

// Stamp computes the hash of the given file and
// creates a timestamp file with it.
func Stamp(logger logging.Logger, filepath, output string, calendars []string) error {
	logger.Info("opening file", logging.String("path", filepath))
	filename := path.Base(filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return err
	}

	digest := hasher.Sum(nil)

	t := serializer.NewTimeStampFile(digest, operation.NewOpSha256())
	for i := range calendars {
		calendar := client.NewCalendar(calendars[i])
		fmt.Printf("Submitting Sha256 digest %s\n to %s",
			hex.EncodeToString(digest), calendar.URL)

		err = calendar.Submit(context.Background(), t.Timestamp, digest)
		if err != nil {
			return err
		}
	}

	if output == "" {
		output = fmt.Sprintf("%s.ots", filename)
	}
	otsFile, err := os.Create(output)
	if err != nil {
		return err
	}

	return serializer.TimestampFileToWriter(t, otsFile)
}
