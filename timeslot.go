package myradio

import (
	"encoding/json"
	"time"
)

type CurrentAndNext struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

type Show struct {
	Title        string `json:"title"`
	Desc         string `json:"desc"`
	Photo        string `json:"photo"`
	StartTimeRaw int64 `json:"start_time"`
	StartTime    time.Time
	EndTimeRaw   int64 `json:"end_time"`
	EndTime      time.Time
	Presenters   string `json:"presenters,omitempty"`
	Url          string `json:"url,omitempty"`
	Id           uint64 `json:"id,omitempty"`
}

func (s *Session) GetCurrentAndNext() (*CurrentAndNext, error) {

	data, err := s.apiRequest("/timeslot/currentandnext", []string{})

	if err != nil {
		return nil, err
	}

	var currentAndNext CurrentAndNext

	err = json.Unmarshal(*data, &currentAndNext)

	if err != nil {
		return nil, err
	}

	currentAndNext.Current.StartTime = time.Unix(currentAndNext.Current.StartTimeRaw, 0)
	currentAndNext.Current.EndTime = time.Unix(currentAndNext.Current.EndTimeRaw, 0)
	currentAndNext.Next.StartTime = time.Unix(currentAndNext.Next.StartTimeRaw, 0)
	currentAndNext.Next.EndTime = time.Unix(currentAndNext.Next.EndTimeRaw, 0)

	return &currentAndNext, nil
}
