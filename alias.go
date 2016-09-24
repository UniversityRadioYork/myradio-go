package myradio

import (
	"encoding/json"
)

type Alias struct {
	Id           int `json:"alias_id"`
	Source       string
	Destinations []struct {
		Atype string `json:"type"`
		Value *json.RawMessage
	}
}

func (s *Session) GetAllAliases(mixins []string) ([]Alias, error) {
	data, err := s.apiRequest("/alias/allaliases", mixins)
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
