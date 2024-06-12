package myradio

import (
	"time"
)

// Term represents information about a MyRadio scheduling term
type Term struct {
	TermID int `json:"term_id"`
	Start int64 `json:"start"`
	Description string `json:"descr"`
	NumWeeks int `json:"num_weeks"`
	WeekNames []string `json:"week_names"`
}

// StartTime returns the start of the term as a time.Time object
func (t *Term) StartTime() (time.Time) {
	return time.Unix(t.Start, 0)
}

// GetAllTerms retrieves all the terms MyRadio is aware of (past and future)
func (s *Session) GetAllTerms() (terms []Term, err error) {
	err = s.get("/term/allterms/").Into(&terms)
	return
}

