package myradio

import (
	"encoding/json"
)

type Banner struct {
	BannerID int    `json:"banner_id"`
	Alt      string `json:"alt"`
	Target   string `json:"Target"`
	URL      string `json:"url"`
}

func (s *Session) GetLiveBanners() (banners []Banner, err error) {
	data, err := s.apiRequest("/banner/livebanners/", []string{})
	if err != nil {
		return
	}

	err = json.Unmarshal(*data, &banners)

	return

}
