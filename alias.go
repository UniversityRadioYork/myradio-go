package myradio

import (
	"encoding/json"

	"github.com/UniversityRadioYork/myradio-go/api"
)

// Alias represents a mail alias.
type Alias struct {
	ID           int `json:"alias_id"`
	Source       string
	Destinations []struct {
		Atype string `json:"type"`
		Value *json.RawMessage
	}
}

// GetAllAliases retrieves all aliases in use.
// It takes a list of additional MyRadio API mixins to use when retrieving the aliases.
// This consumes one API request.
func (s *Session) GetAllAliases(mixins []string) (aliases []Alias, err error) {
	rq := api.NewRequest("/alias/allaliases")
	rq.Mixins = mixins
	err = s.do(rq).Into(&aliases)
	return
}
