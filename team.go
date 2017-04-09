package myradio

import "time"

// Position represents a MyRadio officer position.
type Position struct {
	Team        Team   `json:"team"`
	OfficerID   uint   `json:"officerid"`
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Ordering    uint   `json:"ordering"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

// Officer represents information about an officership inside a Team.
type Officer struct {
	User            User `json:"user"`
	From            time.Time
	FromRaw         int64    `json:"from"`
	MemberOfficerID uint     `json:"memberofficerid"`
	Position        Position `json:"position"`
}

// Team represents a station committee team.
type Team struct {
	TeamID      uint      `json:"teamid"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias"`
	Ordering    uint      `json:"ordering"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Officers    []Officer `json:"officers"`
}

// HeadPosition represents the head position of a team.
type HeadPosition struct {
	User            User
	From            int
	MemberOfficerID int
	Position        OfficerPosition
}

// GetCurrentTeams retrieves all teams inside the station committee.
// This consumes one API request.
func (s *Session) GetCurrentTeams() (teams []Team, err error) {
	err = s.get("/team/currentteams/").into(&teams)
	return
}

// GetTeamWithOfficers retrieves a team record with officer information for the given team name.
// This consumes one API request.
func (s *Session) GetTeamWithOfficers(teamName string) (team Team, err error) {
	if err = s.getf("/team/byalias/%s", teamName).mixin("officers").into(&team); err != nil {
		return
	}

	for k, v := range team.Officers {
		team.Officers[k].From = time.Unix(v.FromRaw, 0)
	}

	return
}

// GetTeamHeadPositions retrieves all head-of-team positions for a given team ID.
// The amount of detail can be controlled using MyRadio mixins.
// This consumes one API request.
func (s *Session) GetTeamHeadPositions(id int, mixins []string) (head []HeadPosition, err error) {
	if err = s.getf("/team/%d/headpositions").mixin(mixins...).into(&head); err != nil {
		return
	}

	for k, v := range head {
		if v.Position.History != nil {
			for ik, iv := range v.Position.History {
				head[k].Position.History[ik].From = time.Unix(iv.FromRaw, 0)
				head[k].Position.History[ik].To = time.Unix(iv.ToRaw, 0)
			}
		}
	}

	return
}
