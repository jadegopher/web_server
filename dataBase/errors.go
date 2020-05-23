package dataBase

import "errors"

var (
	UserNotFoundError   = errors.New("100 user not found")
	NicknameUniqueError = errors.New("101 'userId' not unique")
	EmailUniqueError    = errors.New("102 'email' already used")
	FieldTooLongError   = errors.New("103 field too long (max length ")
	WrongValueError     = errors.New("104 wrong value for field ")
	WrongSymbolsError   = errors.New("105 wrong symbols")
	LoginError          = errors.New("106 login error")
)
