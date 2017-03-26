package myradio

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// parseShortTime parses times in MyRadio's 'DD/MM/YYYY HH:MM' local-time format.
// On success, it returns the equivalent time; else, it reports an error.
func parseShortTime(value string) (time.Time, error) {
	return time.ParseInLocation("02/01/2006 15:04", value, time.Local)
}

// parseDuration parses durations in MyRadio's 'HH:MM:SS' format.
// On success, it returns the equivalent duration; else, it reports an error.
// Errors occur if the duration is not in the given format, or cannot be represented as a Duration.
// For a negative duration, give (-HH:MM:SS).
func parseDuration(value string) (dur time.Duration, err error) {
	// @MattWindsor91:
	// Previously, we relied on time.Parse to implement this.
	// However, while convenient, this fails for durations that don't 'look' like times,
	// eg. anything over 23:59:60.

	vs := strings.Split(value, ":")
	if len(vs) != 3 {
		err = fmt.Errorf("parseDuration: '%s' has %i sections but should have 3", value, len(vs))
		return
	}

	// If the entire thing is prefixed by a sign, then we need to handle that carefully.
	// This is because just treating the hours as negative fails if the hours part is -0!
	sign := 1
	if 0 < len(vs[0]) && vs[0][0] == '-' {
		vs[0] = strings.Replace(vs[0], "-", "0", -1)
		sign = -1
	}

	// Save us from repeating the same processing logic 3 times, at the cost of some efficiency when we fail.
	parseBit := func(bit string, min, max int64) (val int64) {
		if err != nil {
			return
		}
		val, err = strconv.ParseInt(bit, 10, 64)
		if err != nil {
			return
		}

		// At this stage, val will contain the parsed value: this is just an additional sanity check.
		if val < min || max < val {
			err = fmt.Errorf("parseDuration: expected %i-%i, got %i", min, max, val)
		}

		return
	}

	h := parseBit(vs[0], math.MinInt64, math.MaxInt64)
	m := parseBit(vs[1], 0, 59)
	s := parseBit(vs[2], 0, 59) // This is a duration, so we don't need to worry about leap seconds.
	if err != nil {
		return
	}

	dur = time.Duration(sign) * ((time.Duration(h) * time.Hour) + (time.Duration(m) * time.Minute) + (time.Duration(s) * time.Second))
	return
}
