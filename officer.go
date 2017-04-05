package myradio

import "time"

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
	Current     []User `json:"current,omitempty"`
	History     []struct {
		User            User
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
func (s *Session) GetAllOfficerPositions(mixins []string) (positions []OfficerPosition, err error) {
	if err = s.apiRequestInto(&positions, "/officer/allofficerpositions", mixins); err != nil {
		return
	}

	for k, v := range positions {
		for ik, iv := range v.History {
			positions[k].History[ik].From = time.Unix(iv.FromRaw, 0)
			positions[k].History[ik].To = time.Unix(iv.ToRaw, 0)
		}
	}

	return
}
