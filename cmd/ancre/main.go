package main

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "info":
		err = infos(os.Args[2])
	case "stamp":
		err = stamp(os.Args[2])
	case "help":
		fmt.Println("list of subcommands: infos, stamp")
	}
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}

func infos(filepath string) error {
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

func stamp(filepath string) error {
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
