package entities

import "database/sql"

type Task struct {
	Name              string         `json:"name"`
	Picture           sql.NullString `json:"picture"`
	BackgroundPicture sql.NullString `json:"backgroundPicture"`
	Description       string         `json:"description"`
	RecommendedTime   string         `json:"recommendedTime"`
}
