package myradio

import (
	"errors"
	"fmt"
	"time"
)

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

// Return a schedule
func (s *Session) PadWithJukebox(schedule Schedule) (err error) {
	jbShow, err := s.GetShow(1)
	if err != nil {
		return
	}
	jbSeason := Season{ShowMeta: *jbShow}

	startTimes, endTimes, err := GetScheduleDayLimits(schedule)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range schedule {
		schedule[k], err = padDayWithJukebox(v, startTimes[k], endTimes[k], jbSeason)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

// Return a day padded with jukebox Timeslots, with start and end limits (full date and time requried)
func padDayWithJukebox(day Day, from time.Time, to time.Time, jbSeason Season) (retday Day, err error) {
	for i := 0; i < len(day); i++ {
		d := day[i].StartTime.Sub(from)
		if d < 0 {
			err = errors.New("Panic!")
			return
		}
		if d > 0 {
			jbTimeslot := Timeslot{Season: jbSeason, StartTime: from, Duration: d}
			retday = append(retday, jbTimeslot)
		}
		retday = append(retday, day[i])
		from = day[i].EndTime()
	}
	if from.Before(to) {
		jbTimeslot := Timeslot{Season: jbSeason, StartTime: from, Duration: to.Sub(from)}
		retday = append(retday, jbTimeslot)
	}
	return
}

// First day must be schedule["1"]
func GetScheduleDayLimits(schedule Schedule) (startTimes []time.Time, endTimes []time.Time, err error) {
	if len(schedule) == 0 {
		err = errors.New("schedule invalid")
	}
	// Initial value
	day := schedule[0]
	startOffset := day[0].StartTime.Sub(midnight(day[0].StartTime))
	endOffset := day[len(day)-1].EndTime().Sub(midnight(day[len(day)-1].StartTime))
	// Find the extremities of times within the schedule
	for i := 0; i < len(schedule); i++ {
		day = schedule[i]
		st := day[0].StartTime
		mid := midnight(st)
		et := day[len(day)-1].EndTime()

		so := st.Sub(mid)
		eo := et.Sub(mid)
		if so < startOffset {
			startOffset = so
		}
		if eo > endOffset {
			endOffset = eo
		}
	}
	// Generate times for output
	startTimes = make([]time.Time, len(schedule))
	endTimes = make([]time.Time, len(schedule))
	for i := 0; i < len(schedule); i++ {
		day = schedule[i]
		mid := midnight(day[0].StartTime)
		startTimes[i] = mid.Add(startOffset)
		endTimes[i] = mid.Add(endOffset)

	}
	return
}

// Return midnight on t
func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
