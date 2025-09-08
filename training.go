package myradio

import (
	"time"

	"github.com/UniversityRadioYork/myradio-go/api"
)

type TrainingSession struct {
	DemoID            string `json:"demo_id"`
	PresenterStatusID string `json:"presenterstatusid"`
	StartTimeRaw      string `json:"demo_time"`
	Host              string `json:"member"`
	HostMemberID      int    `json:"memberid"`
}

func (ts *TrainingSession) StartTime() time.Time {
	t, _ := time.Parse("Mon 02 Jan 15:04", ts.StartTimeRaw)

	return time.Date(time.Now().Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
}

func (s *Session) GetFutureTrainingSessions() (sessions []TrainingSession, err error) {
	rq := api.NewRequestf("/demo/listdemos")
	err = s.do(rq).Into(&sessions)

	return
}
