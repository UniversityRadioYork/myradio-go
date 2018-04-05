package myradio

const (
	// Values for the current selection/where it was selected from
	SelectorStudio1 = 1
	SelectorStudio2 = 2
	SelectorJukebox = 3
	SelectorOffAir  = 8
	SelectorOB      = 4
	SelectorAux     = 0
	SelectorHub     = 3
	// Values for the selector lock
	LockOff = 0
	LockAux = 1
	LockKey = 2
	// Values for studio power
	OnNone = 0
	OnS1   = 1
	OnS2   = 2
	OnBoth = 3
)

// SelectorInfo holds data from the /selector/query endpoint
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
	return
}
