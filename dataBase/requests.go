package dataBase

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"time"
	"web_server/entities"
)

type DataBase struct {
	Connection *sql.DB
}

type DBConfig struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	DbName         string `json:"dbName"`
	Ip             string `json:"ip"`
	Port           int    `json:"port"`
	ReInitDataBase bool   `json:"reInitDataBase"`
}

func NewDataBase(config *DBConfig) (*DataBase, error) {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Ip, config.Port, config.Username, config.Password, config.DbName)
	conn, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	ret := &DataBase{
		Connection: conn,
	}
	if config.ReInitDataBase {
		if err = ret.createNewDataBase(); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (db *DataBase) AddUser(userInfo *entities.Registration) error {
	if err := db.validateUserInfo(userInfo); err != nil {
		return err
	}
	if err := db.append("INSERT INTO user_private VALUES($1, $2, $3)",
		userInfo.UserId, userInfo.Email, userInfo.Password); err != nil {
		return err
	}
	if err := db.append("INSERT INTO user_info VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
		userInfo.UserId, userInfo.FirstName, userInfo.LastName,
		userInfo.RegistrationTime.Format(time.RFC1123), userInfo.Gender,
		userInfo.OnlineTime.Format(time.RFC1123), userInfo.Picture, userInfo.BackgroundPicture); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) Login(private *entities.UserPrivate) (string, error) {
	result, err := db.Connection.Query(`SELECT * 
		FROM user_private
 		WHERE (user_id = $1 OR email = $2) 
 		AND password = $3`, private.UserId, private.Email, private.Password)
	if err != nil {
		return "", err
	}
	ret := &entities.UserPrivate{}
	if result.Next() {
		if err := result.Scan(&ret.UserId, &ret.Email, &ret.Password); err != nil {
			return "", err
		}
	} else {
		return "", LoginError
	}
	//TODO Update Online INfo
	return ret.UserId, nil
}

func (db *DataBase) GetUserInfo(id string) (*entities.UserInfo, error) {
	result, err := db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	var ret *entities.UserInfo
	if result.Next() {
		ret, err = db.getUserInfo(result)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, UserNotFoundError
	}
	return ret, nil
}

func (db *DataBase) SearchUser(query, from, count string) ([]entities.UserInfo, error) {
	result, err := db.Connection.Query(`SELECT *
		FROM user_info
		WHERE LOWER(user_id) LIKE LOWER($1)
		OR LOWER(first_name) LIKE LOWER($1)
		OR LOWER(last_name) LIKE LOWER($1) LIMIT $2 OFFSET $3`, query+"%", count, from)
	if err != nil {
		return nil, err
	}
	ret := make([]entities.UserInfo, 0)
	for result.Next() {
		elem, err := db.getUserInfo(result)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *elem)
	}
	return ret, nil
}

func (db *DataBase) LogAdd(logInfo *entities.Log) error {
	if err := db.append("INSERT INTO log VALUES($1, $2, $3, $4, $5, $6)",
		logInfo.Time, logInfo.Request, logInfo.Error, logInfo.Body, logInfo.Query,
		logInfo.Headers); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) AddDeveloper(userId string) error {
	if err := db.append("INSERT into developers VALUES($1)", userId); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) CheckDeveloper(userId string) error {
	result, err := db.Connection.Query("SELECT * FROM developers WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	if !result.Next() {
		return UserNotFoundError
	}
	return nil
}

func (db *DataBase) AddTag(tag *entities.Tag) error {
	if err := db.validateTag(tag); err != nil {
		return err
	}
	if err := db.append("INSERT into tags VALUES($1, $2)", tag.Name, tag.Description); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) AddTask(task *entities.Task) error {
	if err := db.validateTask(task); err != nil {
		return err
	}
	if err := db.append("INSERT into tasks VALUES($1, $2, $3, $4, $5)",
		task.Name, task.Description, task.RecommendedTime,
		task.Picture, task.BackgroundPicture); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) validateUserInfo(userInfo *entities.Registration) error {
	r := regexp.MustCompile(nickReg)
	if !r.MatchString(userInfo.UserId) {
		return WrongSymbolsError
	}
	if len(userInfo.UserId) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", userIdField)
	}
	if len(userInfo.Email) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", emailField)
	}
	if len(userInfo.Password) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", passwordField)
	}
	if len(userInfo.FirstName) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", firstNameField)
	}
	if len(userInfo.LastName) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", lastNameField)
	}
	if userInfo.Gender != male && userInfo.Gender != female && userInfo.Gender != another {
		return db.errorConstructValue(WrongValueError, "gender", male, female, another)
	}
	if len(userInfo.Picture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", pictureField)
	}
	if len(userInfo.BackgroundPicture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", bgPictureField)
	}
	result, err := db.Connection.Exec("SELECT user_id FROM user_private WHERE user_id = $1", userInfo.UserId)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 1 {
		return NicknameUniqueError
	}
	result, err = db.Connection.Exec("SELECT email FROM user_private WHERE email = $1", userInfo.Email)
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
		return db.errorConstructLong(FieldTooLongError, "50", "tag_name")
	}
	if len(tag.Description) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", "description")
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
		return db.errorConstructLong(FieldTooLongError, "50", "task_name")
	}
	if len(task.Description) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", "description")
	}
	if len(task.Picture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", "picture")
	}
	if len(task.BackgroundPicture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", "backgroundPicture")
	}
	if len(task.RecommendedTime) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", "recommendedTime")
	}
	return nil
}
