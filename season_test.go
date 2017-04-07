package myradio

// Tests for internal methods of season.go.
// For tests of public methods, see public_tests/season_test.go

import "testing"

// TestIsScheduled tests the isScheduled private method of the Season struct.
func TestIsScheduled(t *testing.T) {
	// TODO(MattWindsor91): this might be a little too low-level
	cases := []struct{
		expected bool
		ftr      string
	}{
		{true, "02/01/2006 15:04"},
		{false, "Not Scheduled"},
	}

	for _, c := range cases {
		s := Season{}
		s.FirstTimeRaw = c.ftr
		got := s.isScheduled()
		if c.expected != got {
			t.Error("with FTR:", c.ftr, "expected:", c.expected, "got:", got)
		}
	}
}
