package myradio

import (
	"encoding/json"
	"fmt"
	"time"
)

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

type Officer struct {
	User            Member `json:"user"`
	From            time.Time
	FromRaw         int64    `json:"from"`
	MemberOfficerID uint     `json:"memberofficerid"`
	Position        Position `json:"position"`
}

type Team struct {
	TeamID      uint      `json:"teamid"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias"`
	Ordering    uint      `json:"ordering"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Officers    []Officer `json:"officers"`
}

type HeadPosition struct {
	User            Member
	From            int
	MemberOfficerID int
	Position        OfficerPosition
}

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

func (s *Session) GetTeamHeadPositions(id int, mixins []string) (head []HeadPosition, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/team/%d/headpositions", id), mixins)
	if err != nil {
		return
	}
	err = json.Unmarshal(*data, &head)
	if err != nil {
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
