package myradio

import (
	"testing"
	"time"
)

func TestParseShortTime(t *testing.T) {
	tests := []struct {
		expected time.Time
		time     string
	}{
		{
			time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local),
			"01/01/1970 00:00",
		},
		{
			time.Date(2009, time.April, 13, 11, 11, 0, 0, time.Local),
			"13/04/2009 11:11",
		},
	}

	for _, test := range tests {
		got, err := parseShortTime(test.time)
		if err != nil {
			t.Error("unexpected error:", err)
		}
		if !got.Equal(test.expected) {
			t.Error("expected:", test.expected, "got:", got)
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		expectedStr string
		time        string
	}{
		{"2h", "02:00:00"},
		{"30m", "00:30:00"},
		{"30h", "30:00:00"},
		{"-5s", "-0:00:05"},
	}

	for _, test := range tests {
		// We can safely leave testing the time lib to the stl
		expected, _ := time.ParseDuration(test.expectedStr)
		got, err := parseDuration(test.time)
		if err != nil {
			t.Error("unexpected error:", err)
		}
		if got != expected {
			t.Error("expected:", expected, "got:", got)
		}
	}
}

func TestParseDurationError(t *testing.T) {
	tests := []struct {
		errStr string
		time   string
	}{
		{"parseDuration: duration string empty", ""},
		{"parseDuration: '01' has 1 sections but should have 3", "01"},
		{"parseDuration: '02:05' has 2 sections but should have 3", "02:05"},
		{"parseDuration: '01:02:03:04' has 4 sections but should have 3", "01:02:03:04"},
		{"parseDuration: expected 0-59, got -1", "10:-1:00"},
		{"parseDuration: expected 0-59, got 60", "10:10:60"},
		{`strconv.ParseInt: parsing "a": invalid syntax`, "a:b:c"},
	}

	for _, test := range tests {
		_, err := parseDuration(test.time)
		if err == nil {
			t.Error("no error, was expecting one")
		}
		if err.Error() != test.errStr {
			t.Error("expected:", test.errStr, "got:", err.Error())
		}
	}
}
