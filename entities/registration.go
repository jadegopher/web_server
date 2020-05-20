package entities

import "time"

type Registration struct {
	Email             string    `json:"email"`
	UserId            string    `json:"userId"`
	Password          string    `json:"password"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	RegistrationTime  time.Time `json:"regTime"`
	Picture           string    `json:"picture"`
	BackgroundPicture string    `json:"backgroundPicture"`
	Gender            string    `json:"gender"`
	OnlineTime        time.Time `json:"onlineTime"`
}
