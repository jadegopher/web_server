package entities

type UserTags struct {
	UserID  string `json:"userId"`
	TagName string `json:"tagName"`
	Rating  int    `json:"rating"`
}
