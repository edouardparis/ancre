package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/ulule/ancre/timestamp"
)

// Remote calendar is an opentimestamp remote calendar
type RemoteCalendar struct {
	URL string
}

// Submit submits the timestamp to the calendar
func (rc RemoteCalendar) Submit(ctx context.Context, t *timestamp.Timestamp, digest []byte) error {
	client := &http.Client{}
	url := fmt.Sprintf("%s/digest", rc.URL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(digest))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to submit to %s", rc.URL)
	}

	return t.Decode(ctx, resp.Body, digest)
}

// NewRemoteCalendar returns a new remote calendar.
func NewRemoteCalendar(url string) *RemoteCalendar {
	return &RemoteCalendar{URL: url}
}
