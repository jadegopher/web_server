package entities

import (
	"time"
)

type UserInfo struct {
	UserId            string    `json:"userId"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	RegistrationTime  time.Time `json:"regTime"`
	Private           int       `json:"private"`
	Picture           string    `json:"picture"`
	BackgroundPicture string    `json:"backgroundPicture"`
	Gender            string    `json:"gender"`
	OnlineTime        time.Time `json:"onlineTime"`
}
