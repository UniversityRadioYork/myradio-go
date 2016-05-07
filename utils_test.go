package myradio

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	expected, err := time.ParseDuration("2h")
	if err != nil {
		t.Fail()
	}
	got, err := parseDuration("15:04:05", "02:00:00")
	if got != expected {
		t.Log("Got:", got, ", Expected:", expected)
		t.Fail()
	}
	expected, err = time.ParseDuration("30m")
	if err != nil {
		t.Fail()
	}
	got, err = parseDuration("15:04:05", "00:30:00")
	if got != expected {
		t.Log("Got:", got, ", Expected:", expected)
		t.Fail()
	}
}
