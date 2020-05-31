package entities

type Quest struct {
	QuestId      int    `json:"questId"`
	UserId       string `json:"userId"`
	TaskName     string `json:"taskName"`
	UserOpponent string `json:"userOpponent"`
	Status       int    `json:"status"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	DeadlineTime string `json:"deadlineTime"`
}
