package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// LogEntry represents a sloppy log entry.
type LogEntry struct {
	Project   *string    `json:"project,omitempty"`
	Service   *string    `json:"service,omitempty"`
	App       *string    `json:"app,omitempty"`
	CreatedAt *Timestamp `json:"createdAt,omitempty"`
	Log       *string    `json:"body,omitempty"`
}

// String prints a log entry
func (e *LogEntry) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		e.CreatedAt.Format("2006-01-02 15:04:05"), *e.Project, *e.Service, *e.App, *e.Log)
}

// Timestamp represents a sloppy timestamp.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON decodes sloppy's date format.
func (u *Timestamp) UnmarshalJSON(data []byte) error {
	var aux int
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("timestamp not a number, got %s", string(data))
	}
	u.Time = time.Unix(int64(aux/1000), 0).UTC()
	return nil
}

func retrieveLogs(c *Client, urlStr string, limit int) (*LogEntry, error) {
	req, err := c.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	// Add limit parameter
	if limit > 0 {
		values := req.URL.Query()
		values.Add("lines", strconv.Itoa(limit))
		req.URL.RawQuery = values.Encode()
	}

	var log LogEntry
	_, err = c.Do(req, log)
	if err != nil {
		return nil, err
	}

	return &log, err
}
