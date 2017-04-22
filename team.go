package myradio

import (
	"encoding/json"
	"fmt"
	"time"
)

// Officer represents information about an officership inside a Team.
type Officer struct {
	User            User `json:"user"`
	From            time.Time
	FromRaw         int64           `json:"from"`
	MemberOfficerID uint            `json:"memberofficerid"`
	Position        OfficerPosition `json:"position"`
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

// GetCurrentTeams retrieves all teams inside the station committee.
// This consumes one API request.
func (s *Session) GetCurrentTeams() (teams []Team, err error) {
	data, err := s.apiRequest("/team/currentteams/", []string{})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &teams)
	if err != nil {
		return
	}
	return
}

// GetTeamWithOfficers retrieves a team record with officer information for the given team name.
// This consumes one API request.
func (s *Session) GetTeamWithOfficers(teamName string) (team Team, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/team/byalias/%s", teamName), []string{"officers"})
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &team)
	if err != nil {
		return
	}
	for k, v := range team.Officers {
		team.Officers[k].From = time.Unix(v.FromRaw, 0)
	}
	return
}

// getTeamPositions retrieves the positions for a given team ID and position type.
// The amount of detail can be controlled using MyRadio mixins.
// The position parameterType is either assistant or head
// This consumes one API request.
func getTeamPositions(positionType string, id int, mixins []string, s *Session) (position []Officer, err error) {
	if positionType != "assistanthead" && positionType != "head" {
		return
	}
	data, err := s.apiRequest(fmt.Sprintf("/team/%d/%spositions", id, positionType), mixins)
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &position)
	if err != nil {
		return
	}
	for k, v := range position {
		position[k].From = time.Unix(v.FromRaw, 0)
		if v.Position.History != nil {
			for ik, iv := range v.Position.History {
				position[k].Position.History[ik].From = time.Unix(iv.FromRaw, 0)
				position[k].Position.History[ik].To = time.Unix(iv.ToRaw, 0)
			}
		}
	}
	return
}

// GetTeamHeadPositions retrieves all head-of-team positions for a given team ID.
// The amount of detail can be controlled using MyRadio mixins.
// This consumes one API request.
func (s *Session) GetTeamHeadPositions(id int, mixins []string) (head []Officer, err error) {
	return getTeamPositions("head", id, mixins, s)
}

// GetTeamAssistantHeadPositions retrieves all assistant-head-of-team positions for a given team ID.
// The amount of detail can be controlled using MyRadio mixins.
// This consumes one API request.
func (s *Session) GetTeamAssistantHeadPositions(id int, mixins []string) (assHead []Officer, err error) {
	return getTeamPositions("assistanthead", id, mixins, s)

}
