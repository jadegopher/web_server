package dataBase

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"regexp"
	"time"
	"web_server/entities"
)

type DataBase struct {
	Connection *sql.DB
}

func NewDataBase(username, password, dbName, ip string, port int) (*DataBase, error) {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		ip, port, username, password, dbName)
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

//TODO return userId
func (db *DataBase) Login(private *entities.UserPrivate) error {
	result, err := db.Connection.Exec("SELECT * FROM user_private WHERE (user_id = $1 OR email = $2) AND password = $3",
		private.UserId, private.Email, private.Password)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 1 {
		return LoginError
	}
	//TODO Update Online INfo
	return nil
}

func (db *DataBase) GetUserInfo(id string) (*entities.UserInfo, error) {
	result, err := db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1", id)
	if err != nil {
		return nil, err
	}
	ret := &entities.UserInfo{}
	var regTime, onlineTime string
	if result.Next() {
		if err := result.Scan(&ret.UserId, &ret.FirstName, &ret.LastName, &regTime,
			&ret.Gender, &onlineTime, &ret.Picture, &ret.BackgroundPicture); err != nil {
			return nil, err
		}
	} else {
		return nil, UserNotFoundError
	}
	ret.RegistrationTime, err = time.Parse(time.RFC1123, regTime)
	if err != nil {
		return nil, err
	}
	ret.OnlineTime, err = time.Parse(time.RFC1123, onlineTime)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *DataBase) append(query string, args ...interface{}) error {
	result, err := db.Connection.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
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
		return db.errorConstructLong(FieldTooLongError, "256", userId)
	}
	if len(userInfo.Email) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", email)
	}
	if len(userInfo.Password) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", password)
	}
	if len(userInfo.FirstName) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", firstName)
	}
	if len(userInfo.LastName) > 256 {
		return db.errorConstructLong(FieldTooLongError, "256", lastName)
	}
	if userInfo.Gender != male && userInfo.Gender != female && userInfo.Gender != another {
		return db.errorConstructValue(WrongValueError, "gender", male, female, another)
	}
	if len(userInfo.Picture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", picture)
	}
	if len(userInfo.BackgroundPicture) > 512 {
		return db.errorConstructLong(FieldTooLongError, "512", backgroundPicture)
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

func (db *DataBase) errorConstructLong(err error, size, name string) error {
	return errors.New(err.Error() + size + " bytes for field '" + name + "')")
}

func (db *DataBase) errorConstructValue(err error, name string, values ...string) error {
	var res string
	for _, elem := range values {
		res += elem + "/"
	}
	return errors.New(err.Error() + "'" + name + "'" + " allowed values: " + res)
}
