package myradio

import "time"

// parseDuration takes a custom layout and a value and returns a time.Duration
//
// Not guaranteed to work so be careful with what you pass in.
func parseDuration(layout, value string) (dur time.Duration, err error) {
	// There is probably a more efficient way of doing this, but time.Unix(0,0) didn't want to work
	midnight, err := time.Parse("15:04:05", "00:00:00")
	if err != nil {
		return
	}
	t, err := time.Parse(layout, value)
	if err != nil {
		return
	}
	return t.Sub(midnight), nil
}
