package entities

type Task struct {
	Name              string `json:"name"`
	Picture           string `json:"picture"`
	BackgroundPicture string `json:"backgroundPicture"`
	Description       string `json:"description"`
	RecommendedTime   string `json:"recommendedTime"`
}
