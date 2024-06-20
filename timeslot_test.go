package myradio_test

import (
	"reflect"
	"testing"
	"time"

	myradio "github.com/UniversityRadioYork/myradio-go"
)

// TestCanEntryZero tests whether a zero-valued CurrentAndNext entry returns true for IsZero.
// It does NOT (yet) test the converse.
func TestCanEntryZero(t *testing.T) {
	s := myradio.Show{}
	if !s.StartTime.IsZero() || !s.EndTime.IsZero() {
		t.Error("zero show returns false for IsZero")
	}
}

// TestCanEntryEnds tests whether a CurrentAndNext entry returns something sensible for Ends.
func TestCanEntryEnds(t *testing.T) {
	cases := []struct {
		t myradio.Time
		e bool
	}{
		{t: myradio.Time{}, e: false},
		{t: myradio.Time{time.Date(2009, time.April, 13, 11, 11, 11, 0, time.UTC)}, e: true},
	}

	for _, c := range cases {
		s := myradio.Show{}
		s.EndTime = c.t
		if s.Ends() != c.e {
			t.Error("show with end time", c.t, "gave incorrect Ends() of", s.Ends())
		}
	}
}

// TestGetWeekScheduleZero tests whether GetWeekSchedule handles empty schedules correctly.
func TestGetWeekScheduleZero(t *testing.T) {
	expected := map[int][]myradio.Timeslot{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
		7: {},
	}

	zeroes := [][]byte{[]byte("[]"), []byte("{}")}
	for _, zero := range zeroes {
		session := myradio.StaticMockSession(zero)

		schedule, err := session.GetWeekSchedule(0, 1)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(schedule, expected) {
			t.Error("expected:", expected, "got:", schedule)
		}
	}
}
