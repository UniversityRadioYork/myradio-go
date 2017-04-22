package myradio

// Banner represents the key information about banners
type Banner struct {
	BannerID int    `json:"banner_id"`
	Alt      string `json:"alt"`
	Target   string `json:"target"`
	URL      string `json:"url"`
}

// GetLiveBanners gets the current live banners
// and returns a slice of banners
func (s *Session) GetLiveBanners() (banners []Banner, err error) {
	err = s.apiRequestInto(&banners, "/banner/livebanners/", []string{})
	return
}
