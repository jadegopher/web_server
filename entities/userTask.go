package entities

import "time"

type UserTask struct {
	UserId           string    `json:"userId"`
	UserId1          string    `json:"UserId1"`
	TaskName         string    `json:"taskName"`
	Private          bool      `json:"private"`
	Status           string    `json:"status"`
	UserConfirmTime  time.Time `json:"userConfirm"`
	User1ConfirmTime time.Time `json:"user1Confirm"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
}
