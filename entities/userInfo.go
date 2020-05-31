package entities

import (
	"database/sql"
	"time"
)

type UserInfo struct {
	UserId            string         `json:"userId"`
	FirstName         string         `json:"firstName"`
	LastName          string         `json:"lastName"`
	RegistrationTime  time.Time      `json:"regTime"`
	Private           int            `json:"private"`
	Picture           sql.NullString `json:"picture"`
	BackgroundPicture sql.NullString `json:"backgroundPicture"`
	Gender            string         `json:"gender"`
	OnlineTime        time.Time      `json:"onlineTime"`
}
