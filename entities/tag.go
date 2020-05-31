package entities

import "database/sql"

type Tag struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}
