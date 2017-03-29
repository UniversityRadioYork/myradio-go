package myradio_test

import (
	myradio "github.com/UniversityRadioYork/myradio-go"
	"reflect"
	"testing"
)

// TestGetWeekScheduleZeroArray tests whether GetWeekSchedule handles [] correctly.
func TestGetWeekScheduleZeroArray(t *testing.T) {
	testGetWeekScheduleZero(t, []byte("[]"))
}

// TestGetWeekScheduleZeroObject tests whether GetWeekSchedule handles {} correctly.
func TestGetWeekScheduleZeroObject(t *testing.T) {
	testGetWeekScheduleZero(t, []byte("{}"))
}

// testGetWeekScheduleZero tests whether GetWeekSchedule handles empty schedules correctly.
func testGetWeekScheduleZero(t *testing.T, zero []byte) {
	expected := map[int][]myradio.Timeslot{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
		7: {},
	}

	session, err := myradio.MockSession(zero)
	if err != nil {
		t.Error(err)
	}

	schedule, err := session.GetWeekSchedule(0, 1)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(schedule, expected) {
		t.Errorf("expected:\n%v\n\ngot:\n%v", expected, schedule)
	}
}
