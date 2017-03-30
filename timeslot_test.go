package myradio_test

import (
	myradio "github.com/UniversityRadioYork/myradio-go"
	"reflect"
	"testing"
)


// testGetWeekScheduleZero tests whether GetWeekSchedule handles empty schedules correctly.
func testGetWeekScheduleZero(t *testing.T) {
	expected := map[int][]myradio.Timeslot{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
		7: {},
	}

	zeroes := [][]byte{ []byte("[]"), []byte("{}") }
	for _, zero := range zeroes{
		session, err := myradio.MockSession(zero)
		if err != nil {
			t.Error(err)
		}

		schedule, err := session.GetWeekSchedule(0, 1)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(schedule, expected) {
			t.Error("expected:", expected, "got:", schedule)
		}
	}
}
