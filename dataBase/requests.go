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
		userInfo.UserPrivate.UserId, userInfo.UserPrivate.Email, userInfo.UserPrivate.Password); err != nil {
		return err
	}
	if err := db.append("INSERT INTO user_info VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		userInfo.UserPrivate.UserId, userInfo.UserInfo.FirstName, userInfo.UserInfo.LastName,
		userInfo.UserInfo.RegistrationTime.Format(time.RFC1123), userInfo.UserInfo.Gender,
		userInfo.UserInfo.OnlineTime.Format(time.RFC1123), userInfo.UserInfo.Private,
		userInfo.UserInfo.Picture, userInfo.UserInfo.BackgroundPicture); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) Login(private *entities.Login) (string, error) {
	result, err := db.Connection.Query(`SELECT * 
		FROM user_private
 		WHERE (user_id = $1 OR email = $1)`, private.Identity)
	if err != nil {
		return "", err
	}
	ret := &entities.UserPrivate{}
	if result.Next() {
		if err := result.Scan(&ret.UserId, &ret.Email, &ret.Password); err != nil {
			return "", err
		}
		if ret.Password != private.Password {
			return "", LoginError
		}
	} else {
		return "", UserNotFoundError
	}
	//TODO Update Online INfo
	return ret.UserId, nil
}

func (db *DataBase) GetUserInfo(id string, self bool) (*entities.UserInfo, error) {
	var result *sql.Rows
	var err error
	if self {
		result, err = db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1", id)
		if err != nil {
			return nil, err
		}
	} else {
		result, err = db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1 AND private = 0", id)
		if err != nil {
			return nil, err
		}
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
		WHERE user_id != $1 AND private = 0 AND
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
	userTags, err := db.GetUserTags(userId, true)
	if err != nil {
		return err
	}
	reqTags := make(map[string]entities.Tag)
	for _, tag := range tags {
		if _, in := reqTags[tag.Name]; !in {
			reqTags[tag.Name] = tag
		}
	}
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
			return db.errorConstructNotFound(TagNotFoundError, value.Name)
		}
		if err := db.append("INSERT INTO user_tags VALUES($1, $2, $3)",
			userId, value.Name, 5); err != nil {
			return err
		}
	}
	return nil
}

func (db *DataBase) GetUserTags(userId string, self bool) ([]entities.IdTags, error) {
	var result *sql.Rows
	var err error
	if self {
		result, err = db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1", userId)
		if err != nil {
			return nil, err
		}
	} else {
		result, err = db.Connection.Query("SELECT * FROM user_info WHERE user_id = $1 AND private = 0", userId)
		if err != nil {
			return nil, err
		}
	}
	if !result.Next() {
		return nil, UserNotFoundError
	}
	result, err = db.Connection.Query("SELECT * FROM user_tags WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	userTags, err := db.getTagsList(result)
	if err != nil {
		return nil, err
	}
	return userTags, nil
}

func (db *DataBase) GetTaskInfo(taskName string) (*entities.Task, error) {
	result, err := db.Connection.Query("SELECT * FROM tasks WHERE task_name = $1", taskName)
	if err != nil {
		return nil, err
	}
	ret := &entities.Task{}
	if result.Next() {
		ret, err = db.getTaskInfo(result)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, db.errorConstructNotFound(TaskNotFoundError, taskName)
	}
	return ret, nil
}

func (db *DataBase) GetTaskTags(taskName string) ([]entities.IdTags, error) {
	result, err := db.Connection.Query("SELECT * FROM tasks WHERE task_name = $1", taskName)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, db.errorConstructNotFound(TaskNotFoundError, taskName)
	}
	result, err = db.Connection.Query("SELECT * FROM task_tags WHERE task_name = $1", taskName)
	if err != nil {
		return nil, err
	}
	taskTags, err := db.getTagsList(result)
	if err != nil {
		return nil, err
	}
	return taskTags, nil
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

func (db *DataBase) AddTagsToTask(taskName string, tags []entities.Tag) error {
	taskTags, err := db.GetTaskTags(taskName)
	if err != nil {
		return err
	}
	reqTags := make(map[string]entities.Tag)
	for _, tag := range tags {
		if _, in := reqTags[tag.Name]; !in {
			reqTags[tag.Name] = tag
		}
	}
	for _, taskTag := range taskTags {
		if _, in := reqTags[taskTag.TagName]; !in {
			if _, err := db.Connection.Exec(`DELETE FROM task_tags 
					WHERE tag_name = $1 AND task_name = $2`,
				taskTag.TagName, taskTag.Id); err != nil {
				return err
			}
		} else {
			delete(reqTags, taskTag.TagName)
		}
	}
	for _, value := range reqTags {
		result, err := db.Connection.Query("SELECT * FROM tags WHERE tag_name = $1", value.Name)
		if err != nil {
			return err
		}
		if !result.Next() {
			return db.errorConstructNotFound(TagNotFoundError, value.Name)
		}
		if err := db.append("INSERT INTO task_tags VALUES($1, $2, $3)",
			taskName, value.Name, 5); err != nil {
			return err
		}
	}
	return nil
}
