package myradio

import (
	"fmt"
	"strconv"
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

type TrainingSessionForSignup struct {
	TrainingSession
	SignupCutoffHours int `json:"signup_cutoff_hours"`
	MaxParticipants   int `json:"max_participants"`
	AttendeeCount     int `json:"attendee_count"`
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

func (s *Session) GetFutureTrainingSessionsForSignup() (sessions []TrainingSessionForSignup, err error) {
	rq := api.NewRequestf("/demo/listdemosforsignup")
	err = s.do(rq).Into(&sessions)
	return
}

func (s *Session) AddAttendeeToDemo(demoID int, userID int) (result int, err error) {
	formParams := make(map[string][]string)
	formParams["userid"] = []string{strconv.Itoa(userID)}
	rs := s.post(fmt.Sprintf("/demo/%d/addattendee", demoID), formParams)
	err = rs.Into(&result)
	return
}
