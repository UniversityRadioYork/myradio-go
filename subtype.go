package myradio

// ShowSeasonSubtype gives information about an available show subtype
type ShowSeasonSubtype struct {
	SubtypeID   string `json:"id"`
	Name        string `json:"name"`
	Class       string `json:"class"`
	Description string `json:"description"`
}

// GetAllShowSubtypes returns an array of all ShowSeasonSubtypes
func (s *Session) GetAllShowSubtypes() (subtypes []ShowSeasonSubtype, err error) {
	err = s.get("/showSubtype/all").Into(&subtypes)
	return
}

// GetShowSubtypeByClass returns a ShowSeasonSubtype based on a given class
// Can return nil if no subtype found for given class
func (s *Session) GetShowSubtypeByClass(class string) (subtype ShowSeasonSubtype, err error) {
	subtypes, err := s.GetAllShowSubtypes()
	if err != nil {
		return
	}

	for _, subtype := range subtypes {
		if subtype.Class == class {
			return subtype, err
		}
	}
	return
}
