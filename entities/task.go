package entities

import "time"

type Task struct {
	Name              string    `json:"name"`
	Picture           string    `json:"picture"`
	BackgroundPicture string    `json:"backgroundPicture"`
	Description       string    `json:"description"`
	Likes             uint64    `json:"likes"`
	Dislikes          uint64    `json:"dislikes"`
	CreatedBy         string    `json:"createdBy"`
	RecommendedTime   time.Time `json:"RecommendedTime"`
}
