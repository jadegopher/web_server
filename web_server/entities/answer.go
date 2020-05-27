package entities

type Answer struct {
	Element interface{} `json:"element,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

type Error struct {
	Code  int    `json:"code"`
	Error string `json:"message"`
}
