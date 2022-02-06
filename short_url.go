package myradio

import (
	"bytes"
	"net/url"
)

// ShortURL represents the key information of a short URL.
type ShortURL struct {
	ShortURLID uint   `json:"short_url_id"`
	Slug       string `json:"slug"`
	RedirectTo string `json:"redirect_to"`
}

func (s *Session) GetAllShortURLs() (urls []ShortURL, err error) {
	err = s.get("/shortUrl/all").Into(&urls)
	return
}

func (s *Session) LogShortURLClick(id uint, userAgent, ipAddress string) error {
	params := url.Values{}
	params["userAgent"] = []string{userAgent}
	params["ipAddress"] = []string{ipAddress}
	resp := s.putf("/shortUrl/%d/logclick", *bytes.NewBufferString(params.Encode()), id)
	_, err := resp.JSON()
	return err
}
