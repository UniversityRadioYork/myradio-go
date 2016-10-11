package myradio

import (
	"encoding/json"
)

// Banner represents the key information about banners
type Banner struct {
	BannerID int    `json:"banner_id"`
	Alt      string `json:"alt"`
	Target   string `json:"Target"`
	URL      string `json:"url"`
}

// GetLiveBanners gets the current live banners
// and returns a sice of banners
func (s *Session) GetLiveBanners() (banners []Banner, err error) {
	data, err := s.apiRequest("/banner/livebanners/", []string{})
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &banners)

	return

}
