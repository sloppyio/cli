package api_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sloppyio/cli/pkg/api"
)

func TestLogEntryString(t *testing.T) {
	input := &api.LogEntry{
		Project:   api.String("letschat"),
		Service:   api.String("frontend"),
		App:       api.String("apache"),
		CreatedAt: &api.Timestamp{time.Date(2015, 12, 07, 20, 10, 0, 0, time.UTC)},
		Log:       api.String("WARN: 123"),
	}
	want := "2015-12-07 20:10:00 letschat frontend apache WARN: 123"

	if input.String() != want {
		t.Errorf("String(%v) = %s, want %s", input, input.String(), want)
	}
}

func TestTimestampUnmarshal(t *testing.T) {
	input := []byte(`1449519000000`)
	want := time.Date(2015, 12, 07, 20, 10, 0, 0, time.UTC)

	timestamp := api.Timestamp{}
	if err := timestamp.UnmarshalJSON(input); err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}
	if timestamp.String() != want.String() {
		t.Errorf("UnmarshalJSON: %v, want %v", timestamp.String(), want.String())
	}
}

func TestTimestampUnmarshal_invalidTimestamp(t *testing.T) {
	input := []byte(`abcd`)

	timestamp := api.Timestamp{}
	if err := timestamp.UnmarshalJSON(input); err == nil {
		t.Error("Expected error to returned")
	}
}

func testLogOutput(t *testing.T, logs <-chan api.LogEntry, errs <-chan error) {
	for i := 0; i < 5; i++ {
		want := fmt.Sprintf("frontend-%d", i)
		select {
		case log, ok := <-logs:
			if ok && *log.Service != want {
				t.Errorf("Log.Service = %v, want %v", *log.Service, want)
			}
		case err := <-errs:
			t.Errorf("Unexpected error: %v", err)
		}
	}
}
