package dataBase

import "errors"

var (
	UserNotFoundError     = errors.New("100 user not found")
	NicknameUniqueError   = errors.New("101 'userIdField' not unique")
	EmailUniqueError      = errors.New("102 'emailField' already used")
	FieldTooLongError     = errors.New("103 field too long (max length ")
	WrongValueError       = errors.New("104 wrong value for field ")
	WrongSymbolsError     = errors.New("105 wrong symbols")
	LoginError            = errors.New("106 login error")
	TagUniqueError        = errors.New("107 tag with this name already exist")
	TagNotFoundError      = errors.New("108 tag with name ")
	TaskNotFoundError     = errors.New("109 task with name ")
	TaskUniqueError       = errors.New("110 task with this name already exist")
	FieldNotFoundError    = errors.New("111 field with name ")
	InviteExistsError     = errors.New("112 invite has already existed")
	InviteToYourselfError = errors.New("113 you can't invite yourself")
	QuestStatusError      = errors.New("114 wrong status")
	QuestNotFoundError    = errors.New("115 quest not found error")
	AppendError           = errors.New("116 append error")
	ExpiredTaskError      = errors.New("117 task expired")
	NoTasksOnServerError  = errors.New("118 can't find any task")
)
