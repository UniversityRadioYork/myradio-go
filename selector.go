package myradio

import (
	"github.com/UniversityRadioYork/myradio-go/api"
)

type SelectorInfo struct {
	Studio       int `json:"studio"`
	Lock         int `json:"lock"`
	SelectedFrom int `json:"selectedfrom"`
	Power        int `json:"power"`
}

// GetSelectorInfo retrieves the current status of the selector
// This consumes one API request.
func (s *Session) GetSelectorInfo() (info *SelectorInfo, err error) {
	err = s.get("/selector/query").Into(&info)
}
