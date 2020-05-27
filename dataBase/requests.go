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

func (db *DataBase) SearchUser(userId, query, from, count string) ([]entities.UserInfo, error) {
	result, err := db.Connection.Query(`SELECT *
		FROM user_info
		WHERE user_id != $1 AND
		(LOWER(user_id) LIKE LOWER($2)
		OR LOWER(first_name) LIKE LOWER($2)
		OR LOWER(last_name) LIKE LOWER($2)) LIMIT $3 OFFSET $4`, userId, query+"%", count, from)
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

func (db *DataBase) DeleteUser(userId string) error {
	if _, err := db.Connection.Exec("DELETE FROM user_private WHERE user_id = $1", userId); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) AddTagsToUser(userId string, tags []entities.Tag) error {
	userTags, err := db.GetUsersTags(userId)
	if err != nil {
		return err
	}
	reqTags := make(map[string]entities.Tag)
	for _, tag := range tags {
		if _, in := reqTags[tag.Name]; !in {
			reqTags[tag.Name] = tag
		}
	}
	fmt.Println(userTags)
	for _, userTag := range userTags {
		if _, in := reqTags[userTag.TagName]; !in {
			if _, err := db.Connection.Exec(`DELETE FROM user_tags 
					WHERE tag_name = $1 AND user_id = $2`,
				userTag.TagName, userId); err != nil {
				return err
			}
		} else {
			delete(reqTags, userTag.TagName)
		}
	}
	for _, value := range reqTags {
		result, err := db.Connection.Query("SELECT * FROM tags WHERE tag_name = $1", value.Name)
		if err != nil {
			return err
		}
		if !result.Next() {
			return db.errorConstructTag(value.Name)
		}
		if err := db.append("INSERT INTO user_tags VALUES($1, $2, $3)",
			userId, value.Name, 5); err != nil {
			return err
		}
	}
	return nil
}

func (db *DataBase) GetUsersTags(userId string) ([]entities.UserTags, error) {
	userTags := make([]entities.UserTags, 0)
	result, err := db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, UserNotFoundError
	}
	result, err = db.Connection.Query("SELECT * FROM user_tags WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		tmp := entities.UserTags{}
		if err := result.Scan(&tmp.UserID, &tmp.TagName, &tmp.Rating); err != nil {
			return nil, err
		}
		userTags = append(userTags, tmp)
	}
	return userTags, nil
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
