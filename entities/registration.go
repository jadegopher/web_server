package entities

type Registration struct {
	UserPrivate UserPrivate `json:"userPrivate"`
	UserInfo    UserInfo    `json:"userInfo"`
}
