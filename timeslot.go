package myradio

import (
	"encoding/json"
)

type CurrentAndNext struct {
	Next    Show `json:"next"`
	Current Show `json:"current"`
}

type Show struct {
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Photo      string `json:"photo"`
	StartTime  uint64 `json:"start_time"`
	EndTime    uint64 `json:"end_time"`
	Presenters string `json:"presenters,omitempty"`
	Url        string `json:"url,omitempty"`
	Id         uint64 `json:"id,omitempty"`
}

func (s *Session) GetCurrentAndNext() (*CurrentAndNext, error) {
	data, err := s.apiRequest("/timeslot/currentandnext", []string)
	if err != nil {
		return nil, err
	}
	var currentAndNext CurrentAndNext
	err = json.Unmarshal(*data, &currentAndNext)
	if err != nil {
		return nil, err
	}
	return &currentAndNext, nil
}
