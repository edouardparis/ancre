package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ulule/ancre/client"
	"github.com/ulule/ancre/operation"
	"github.com/ulule/ancre/serializer"
)

// Stamp computes the hash of the given file and
// creates a timestamp file with it.
func Stamp(filepath string) error {
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

	calendar := client.NewRemoteCalendar("https://a.pool.opentimestamps.org")

	fmt.Printf("Submitting Sha256 digest %s\n to %s",
		hex.EncodeToString(digest), calendar.URL)

	otsFile, err := os.Create(fmt.Sprintf("%s.ots", filename))
	if err != nil {
		return err
	}

	t := serializer.NewTimeStampFile(digest, operation.NewOpSha256())
	err = calendar.Submit(context.Background(), t.Timestamp, digest)
	if err != nil {
		return err
	}

	return serializer.TimestampFileToWriter(t, otsFile)
}
