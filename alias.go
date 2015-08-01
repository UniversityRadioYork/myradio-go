package myradio

import (
	"encoding/json"
)

type Alias struct {
	id           int `json:"alias_id"`
	source       string
	destinations []struct {
		atype string `json:"type"`
		value *json.RawMessage
	}
}

func (s *Session) GetAllAliases() ([]Alias, error) {
	data, err := s.apiRequest("/alias/allaliases", nil)
	if err != nil {
		return nil, err
	}
	var aliases []Alias
	err = json.Unmarshal(*data, &aliases)
	if err != nil {
		return nil, err
	}
	return aliases, nil
}
