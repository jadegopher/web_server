package entities

import "database/sql"

type Quest struct {
	QuestId      int            `json:"questId"`
	UserId       string         `json:"userId"`
	TaskName     sql.NullString `json:"taskName"`
	UserOpponent string         `json:"userOpponent"`
	Status       int            `json:"status"`
	StartTime    sql.NullString `json:"startTime"`
	EndTime      sql.NullString `json:"endTime"`
	DeadlineTime sql.NullString `json:"deadlineTime"`
}
