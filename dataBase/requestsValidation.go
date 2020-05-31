package dataBase

import (
	"regexp"
	"web_server/entities"
)

func (db *DataBase) validateUserInfo(userInfo *entities.Registration) error {
	if userInfo.UserPrivate.UserId == "" {
		return db.errorConstructNotFound(FieldNotFoundError, userIdField)
	}
	if userInfo.UserPrivate.Email == "" {
		return db.errorConstructNotFound(FieldNotFoundError, emailField)
	}
	if userInfo.UserInfo.FirstName == "" {
		return db.errorConstructNotFound(FieldNotFoundError, firstNameField)
	}
	if userInfo.UserInfo.LastName == "" {
		return db.errorConstructNotFound(FieldNotFoundError, lastNameField)
	}
	if userInfo.UserPrivate.Password == "" {
		return db.errorConstructNotFound(FieldNotFoundError, passwordField)
	}
	r := regexp.MustCompile(nickReg)
	if !r.MatchString(userInfo.UserPrivate.UserId) {
		return WrongSymbolsError
	}
	if len(userInfo.UserPrivate.UserId) > 256 {
		return db.errorConstructLong("256", userIdField)
	}
	if len(userInfo.UserPrivate.Email) > 256 {
		return db.errorConstructLong("256", emailField)
	}
	if len(userInfo.UserPrivate.Password) > 256 {
		return db.errorConstructLong("256", passwordField)
	}
	if len(userInfo.UserInfo.FirstName) > 256 {
		return db.errorConstructLong("256", firstNameField)
	}
	if len(userInfo.UserInfo.LastName) > 256 {
		return db.errorConstructLong("256", lastNameField)
	}
	if userInfo.UserInfo.Gender != male && userInfo.UserInfo.Gender != female && userInfo.UserInfo.Gender != another {
		return db.errorConstructValue(WrongValueError, "gender", male, female, another)
	}
	if len(userInfo.UserInfo.Picture.String) > 512 {
		return db.errorConstructLong("512", pictureField)
	}
	if len(userInfo.UserInfo.BackgroundPicture.String) > 512 {
		return db.errorConstructLong("512", bgPictureField)
	}
	if userInfo.UserInfo.Private != 0 && userInfo.UserInfo.Private != 1 {
		return db.errorConstructValue(WrongValueError, "private", "0", "1")
	}
	result, err := db.Connection.Exec("SELECT user_id FROM user_private WHERE user_id = $1", userInfo.UserPrivate.UserId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 1 {
		return NicknameUniqueError
	}
	result, err = db.Connection.Exec("SELECT email FROM user_private WHERE email = $1", userInfo.UserPrivate.Email)
	if err != nil {
		return err
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 1 {
		return EmailUniqueError
	}
	return nil
}

func (db *DataBase) validateTag(tag *entities.Tag) error {
	if len(tag.Name) > 50 {
		return db.errorConstructLong("50", "tag_name")
	}
	if len(tag.Description.String) > 256 {
		return db.errorConstructLong("256", "description")
	}
	result, err := db.Connection.Exec("SELECT tag_name FROM tags WHERE tag_name = $1", tag.Name)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 1 {
		return TagUniqueError
	}
	return nil
}

func (db *DataBase) validateTask(task *entities.Task) error {
	if len(task.Name) > 50 {
		return db.errorConstructLong("50", "task_name")
	}
	if len(task.Description) > 256 {
		return db.errorConstructLong("256", "description")
	}
	if len(task.Picture.String) > 512 {
		return db.errorConstructLong("512", "picture")
	}
	if len(task.BackgroundPicture.String) > 512 {
		return db.errorConstructLong("512", "backgroundPicture")
	}
	if len(task.RecommendedTime) > 256 {
		return db.errorConstructLong("256", "recommendedTime")
	}
	result, err := db.Connection.Exec("SELECT task_name FROM tasks WHERE task_name = $1", task.Name)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 1 {
		return TaskUniqueError
	}
	return nil
}
