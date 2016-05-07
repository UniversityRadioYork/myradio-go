package myradio

import (
	"encoding/json"
	"fmt"
	"time"
)

type Officership struct {
	OfficerId   uint   `json:"officerid,string"`
	OfficerName string `json:"officer_name"`
	TeamId      uint   `json:"teamid,string"`
	FromDateRaw string `json:"from_date,omitempty"`
	FromDate    time.Time
	TillDateRaw string `json:"till_date,omitempty"`
	TillDate    time.Time
}

type Photo struct {
	PhotoId      uint   `json:"photoid"`
	DateAddedRaw string `json:"date_added"`
	DateAdded    time.Time
	Format       string `json:"format"`
	Owner        uint   `json:"owner"`
	Url          string `json:"url"`
}

func (s *Session) GetUserBio(id int) (string, error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d/bio/", id), []string{})
	if err != nil {
		return "", err
	}
	var bio string
	err = json.Unmarshal(*data, &bio)
	if err != nil {
		return "", err
	}
	return bio, nil
}

func (s *Session) GetUserName(id int) (string, error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d/name/", id), []string{})
	if err != nil {
		return "", err
	}
	var name string
	err = json.Unmarshal(*data, &name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *Session) GetUserProfilePhoto(id int) (profilephoto Photo, err error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d/profilephoto/", id), []string{})
	if err != nil {
		return profilephoto, err
	}
	err = json.Unmarshal(*data, &profilephoto)
	if err != nil {
		return profilephoto, err
	}
	profilephoto.DateAdded, err = time.Parse("02/01/2006 15:04", profilephoto.DateAddedRaw)
	if err != nil {
		return profilephoto, err
	}
	return profilephoto, nil
}

func (s *Session) GetUserOfficerships(id int) ([]Officership, error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d/officerships/", id), []string{})
	if err != nil {
		return nil, err
	}
	var officerships []Officership
	err = json.Unmarshal(*data, &officerships)
	if err != nil {
		return nil, err
	}
	for k, v := range officerships {
		if officerships[k].FromDateRaw != "" {
			officerships[k].FromDate, err = time.Parse("2006-01-02", v.FromDateRaw)
			if err != nil {
				return nil, err
			}
		}
		if officerships[k].TillDateRaw != "" {
			officerships[k].TillDate, err = time.Parse("2006-01-02", v.FromDateRaw)
			if err != nil {
				return nil, err
			}
		}
	}
	return officerships, nil
}

func (s *Session) GetUserShowCredits(id int) ([]ShowMeta, error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d/shows/", id), []string{})
	if err != nil {
		return nil, err
	}
	var shows []ShowMeta
	err = json.Unmarshal(*data, &shows)
	if err != nil {
		return nil, err
	}
	return shows, nil
}
