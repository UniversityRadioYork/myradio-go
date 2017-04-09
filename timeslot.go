package myradio

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// CurrentAndNext stores a pair of current and next show.
type CurrentAndNext struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

// Show contains a summary of information about a URY schedule timeslot.
type Show struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Photo        string `json:"photo"`
	StartTimeRaw int64  `json:"start_time"`
	StartTime    time.Time
	EndTimeRaw   string `json:"end_time"` // Sometimes "The End of Time"
	EndTime      time.Time
	Presenters   string `json:"presenters,omitempty"`
	Url          string `json:"url,omitempty"`
	Id           uint64 `json:"id,omitempty"`
}

// IsZero determines whether the Show is the zero value.
// This is useful for checking when, for example, there is no next show in a CurrentAndNext.
func (s *Show) IsZero() bool {
	// We could use reflect.DeepEqual(Show{}) here,
	// but an easier way is to check a field that will _never_ be its zero value.
	// Assume that EndTimeRaw will always be a number, or "The End of Time", not "".
	// This assumption will eventually go stale!
	return s.EndTimeRaw == ""
}

// Ends determines whether the Show has a defined end time.
func (s *Show) Ends() bool {
	// populateShowTimes() will define EndTime as zero if there isn't one.
	return s.EndTime.IsZero()
}

// populateShowTimes sets the times for the given Show given their raw values.
func (s *Show) populateShowTimes() error {
	s.StartTime = time.Unix(s.StartTimeRaw, 0)

	// As mentioned above, sometimes EndTimeRaw is "The End of Time".
	// This is a known MyRadio-ism!
	if s.EndTimeRaw == "The End of Time" {
		// Whatever this is sent to should give 'true' for Show.Ends().
		s.EndTime = time.Time{}
		return nil
	}

	timeint, err := strconv.ParseInt(s.EndTimeRaw, 10, 64)
	if err != nil {
		return err
	}

	s.EndTime = time.Unix(timeint, 0)
	return nil
}

// Timeslot contains information about a single timeslot in the URY schedule.
// A timeslot is a single slice of time on the schedule, typically one hour long.
type Timeslot struct {
	Season
	TimeslotID     uint64   `json:"timeslot_id"`
	TimeslotNum    int      `json:"timeslot_num"`
	Tags           []string `json:"tags"`
	Time           time.Time
	TimeRaw        int64 `json:"time"`
	StartTime      time.Time
	StartTimeRaw   string `json:"start_time"`
	Duration       time.Duration
	DurationRaw    string `json:"duration"`
	MixcloudStatus string `json:"mixcloud_status"`
}

// populateTimeslotTimes sets the times for the given Timeslot given their raw values.
func (t *Timeslot) populateTimeslotTimes() (err error) {
	// Remember: a Timeslot is a supertype of Season.
	if err = t.populateSeasonTimes(); err != nil {
		return
	}

	t.Time = time.Unix(t.TimeRaw, 0)

	t.StartTime, err = parseShortTime(t.StartTimeRaw)
	if err != nil {
		return
	}

	t.Duration, err = parseDuration(t.DurationRaw)
	return
}

// TracklistItem represents a single item in a show tracklist.
type TracklistItem struct {
	Track
	Album        Album `json:"album"`
	EditLink     Link  `json:"editlink"`
	DeleteLink   Link  `json:"deletelink"`
	Time         time.Time
	TimeRaw      int64 `json:"time"`
	StartTime    time.Time
	StartTimeRaw string `json:"starttime"`
	AudioLogID   uint   `json:"audiologid"`
}

// GetCurrentAndNext gets the current and next shows at the time of the call.
// This consumes one API request.
func (s *Session) GetCurrentAndNext() (can *CurrentAndNext, err error) {
	if err = s.get("/timeslot/currentandnext").into(&can); err != nil {
		return
	}

	if err = can.Current.populateShowTimes(); err != nil {
		return
	}

	// Sometimes, we only get a Current, not a Next.
	// Don't try populate times on a show that doesn't exist.
	if can.Next.IsZero() {
		return
	}
	err = can.Next.populateShowTimes()

	return
}

// GetWeekSchedule gets the weekly schedule for ISO 8601 week week of year year.
// If such a schedule exists, it returns the result as an map from ISO 8601 weekdays to timeslot slices.
// Thus, 1 maps to Monday's timeslots; 2 to Tuesday; and so on.
// Each slice progresses chronologically from start of URY day to finish of URY day.
// If no such schedule exists, it returns a map of empty slices.
// If an error occurred, this is returned in error, and the timeslot map is undefined.
// This consumes one API request.
func (s *Session) GetWeekSchedule(year, week int) (map[int][]Timeslot, error) {
	// TODO(CaptainHayashi): proper errors
	if year < 0 {
		return nil, fmt.Errorf("year %d is too low", year)
	}
	if week < 1 || 53 < week {
		return nil, fmt.Errorf("week %d is not within the ISO range 1..53", week)
	}

	data, aerr := s.get(fmt.Sprintf("/timeslot/weekschedule/%d", week)).param("year", strconv.Itoa(year)).do()
	if aerr != nil {
		return nil, aerr
	}

	// MyRadio responds in a different way when the schedule is empty, so we need to catch that.
	// See https://github.com/UniversityRadioYork/MyRadio/issues/665 for details.
	if isEmptySchedule(data) {
		return map[int][]Timeslot{
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
			6: {},
			7: {},
		}, nil
	}

	// The timeslots come to us with string keys labelled with the weekday.
	// These timeslots start from "1" (Monday) and go up to "7" (Sunday).
	// Note that this is different from Go's view of the week (0 = Sunday, 1 = Monday)!
	stringyTimeslots := make(map[string][]Timeslot)
	if jerr := json.Unmarshal(*data, &stringyTimeslots); jerr != nil {
		return nil, jerr
	}

	return destringTimeslots(stringyTimeslots)
}

// isEmptySchedule tries to work out, from MyRadio schedule JSON, whether the schedule is empty.
func isEmptySchedule(data json.Marshaler) bool {
	bs, err := data.MarshalJSON()
	if err != nil {
		// The logic later on in GetWeekSchedule should hit this same error, so handle it there.
		return false
	}

	if len(bs) != 2 {
		return false
	}

	// Due to a quirk in MyRadio (well, PHP), the empty schedule can be returned as the empty array '[]',
	// instead of the empty object '{}'.
	if bs[0] == '[' && bs[1] == ']' {
		return true
	}

	if bs[0] == '{' && bs[1] == '}' {
		return true
	}

	return false
}

// destringTimeslots converts a week schedule from string indices to integer indices.
// It takes a map from strings "1"--"7" to day schedules, and returns a map from integers 1--7 to day schedules.
// It returns an error if any of the string indices cannot be converted.
func destringTimeslots(stringyTimeslots map[string][]Timeslot) (map[int][]Timeslot, error) {
	timeslots := make(map[int][]Timeslot)
	for sday, ts := range stringyTimeslots {
		day, derr := strconv.Atoi(sday)
		if derr != nil {
			return nil, derr
		}
		for i := range ts {
			if terr := ts[i].populateTimeslotTimes(); terr != nil {
				return nil, terr
			}
		}
		timeslots[day] = ts
	}

	return timeslots, nil
}

// GetTimeslot retrieves the timeslot with the given ID.
// This consumes one API request.
func (s *Session) GetTimeslot(id int) (timeslot Timeslot, err error) {
	if err = s.get(fmt.Sprintf("/timeslot/%d", id)).into(&timeslot); err != nil {
		return
	}
	err = timeslot.populateTimeslotTimes()
	return
}

// GetTrackListForTimeslot retrieves the tracklist for the timeslot with the given ID.
// This consumes one API request.
func (s *Session) GetTrackListForTimeslot(id int) (tracklist []TracklistItem, err error) {
	if err = s.get(fmt.Sprintf("/tracklistItem/tracklistfortimeslot/%d", id)).into(&tracklist); err != nil {
		return
	}
	for k, v := range tracklist {
		tracklist[k].Time = time.Unix(tracklist[k].TimeRaw, 0)
		tracklist[k].StartTime, err = time.Parse("02/01/2006 15:04:05", v.StartTimeRaw)
		if err != nil {
			return nil, err
		}
	}
	return
}
