package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/ulule/ancre/format/ots"
	"github.com/ulule/ancre/timestamp"
)

// Calendar calendar is an opentimestamp remote calendar
type Calendar struct {
	URL string
}

// Submit submits the timestamp to the calendar
func (cal Calendar) Submit(ctx context.Context, digest []byte) (*timestamp.Timestamp, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/digest", cal.URL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(digest))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to submit to %s", cal.URL)
	}

	return ots.Decode(ctx, resp.Body, digest)
}

// NewCalendar returns a new remote calendar.
func NewCalendar(url string) *Calendar {
	return &Calendar{URL: url}
}
