package dataBase

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

func (db *DataBase) LogAdd(logInfo *Log) error {
	if err := db.append("INSERT INTO log VALUES($1, $2, $3, $4, $5, $6)",
		logInfo.Time, logInfo.Request, logInfo.Error, logInfo.Body, logInfo.Query,
		logInfo.Headers); err != nil {
		return err
	}
	return nil
}
