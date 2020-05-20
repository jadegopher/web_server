package entities

type Answer struct {
	Element interface{} `json:"element"`
	Error   string      `json:"errorMessage"`
}
