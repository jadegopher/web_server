package entities

type UserPrivate struct {
	Email    string `json:"email"`
	UserId   string `json:"userId"`
	Password string `json:"password"`
}
