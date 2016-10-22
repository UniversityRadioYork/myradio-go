package myradio

import (
	"encoding/json"
	"time"
)

// OfficerPosition represents a specific station officer position.
type OfficerPosition struct {
	OfficerID   int
	Name        string
	Alias       string
	Team        Team
	Ordering    int
	Description string
	Status      string
	Type        string
	Current     []Member `json:"current,omitempty"`
	History     []struct {
		User            Member
		From            time.Time
		FromRaw         int64 `json:"from"`
		To              time.Time
		ToRaw           int64 `json:"to"`
		MemberOfficerID int
	} `json:"history,omitempty"`
}

// GetAllOfficerPositions retrieves all officer positions in MyRadio.
// The amount of detail can be controlled by adding MyRadio mixins.
// This consumes one API request.
func (s *Session) GetAllOfficerPositions(mixins []string) ([]OfficerPosition, error) {
	data, err := s.apiRequest("/officer/allofficerpositions", mixins)
	if err != nil {
		return nil, err
	}
	var positions []OfficerPosition
	err = json.Unmarshal(*data, &positions)
	if err != nil {
		return nil, err
	}
	for k, v := range positions {
		for ik, iv := range v.History {
			positions[k].History[ik].From = time.Unix(iv.FromRaw, 0)
			positions[k].History[ik].To = time.Unix(iv.ToRaw, 0)
		}
	}
	return positions, nil
}
