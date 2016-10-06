package myradio

import (
	"encoding/json"
)

// Alias represents a mail alias.
type Alias struct {
	Id           int `json:"alias_id"`
	Source       string
	Destinations []struct {
		Atype string `json:"type"`
		Value *json.RawMessage
	}
}

// GetAllAliases retrieves all aliases in use.
//
// It takes a list of additional MyRadio API mixins to use when retrieving the aliases.
//
// This consumes one API request.
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
