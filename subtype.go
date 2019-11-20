package myradio

type ShowSeasonSubtype struct {
	SubtypeID string `json:"id"`
	Name      string `json:"name"`
	Class     string `json:"class"`
}

func (s *Session) GetAllShowSubtypes() (subtypes []ShowSeasonSubtype, err error) {
	err = s.get("/showSubtype/all").Into(&subtypes)
	return
}
