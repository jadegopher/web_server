package dataBase

import "errors"

var (
	UserNotFoundError   = errors.New("100 user not found")
	NicknameUniqueError = errors.New("101 'userIdField' not unique")
	EmailUniqueError    = errors.New("102 'emailField' already used")
	FieldTooLongError   = errors.New("103 field too long (max length ")
	WrongValueError     = errors.New("104 wrong value for field ")
	WrongSymbolsError   = errors.New("105 wrong symbols")
	LoginError          = errors.New("106 login error")
	TagUniqueError      = errors.New("107 tag with this name already exists")
)
