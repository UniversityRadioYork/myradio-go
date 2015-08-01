package myradio

import (
	"encoding/json"
	"fmt"
)

type Member struct {
	Memberid     int
	Fname, Sname string
	Sex          string
	Email        string `json:"public_email"`
	Receiveemail bool   `json:"receive_email"`
}

func (s *Session) GetMember(id int) (*Member, error) {
	data, err := s.apiRequest(fmt.Sprintf("/user/%d", id), []string{"personal_data"})
	if err != nil {
		return nil, err
	}
	var member Member
	err = json.Unmarshal(*data, &member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
