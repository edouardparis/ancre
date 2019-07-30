package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/edouardparis/ancre/calendar/client"
	"github.com/edouardparis/ancre/format/ots"
	"github.com/edouardparis/ancre/logging"
	"github.com/edouardparis/ancre/operation"
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

	t := ots.NewTimeStampFile(digest, operation.NewOpSha256())
	for i := range calendars {
		calendar := client.NewCalendar(calendars[i])
		logger.Info("Submitting to calendar",
			logging.String("url", calendar.URL),
			logging.String("digest", hex.EncodeToString(digest)))

		ts, err := calendar.Submit(context.Background(), digest)
		if err != nil {
			return err
		}
		err = t.Timestamp.Merge(ts)
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

	return ots.TimestampFileToWriter(t, otsFile)
}
