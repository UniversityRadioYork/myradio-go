package myradio

import (
	"testing"
	"time"
)

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
		if err != nil || got != expected {
			t.Error("Got:", got, ", Expected:", expected, ", Error:", err)
		}
	}
}
