package dataBase

import (
	"database/sql"
	"errors"
	"time"
	"web_server/entities"
)

func (db *DataBase) createNewDataBase() error {
	if err := db.createLogTable(); err != nil {
		return err
	}
	if err := db.createUserPrivateTable(); err != nil {
		return err
	}
	if err := db.createUserInfoTable(); err != nil {
		return err
	}
	if err := db.createDevelopersTable(); err != nil {
		return err
	}
	if err := db.createTagsTable(); err != nil {
		return err
	}
	if err := db.createTasksTable(); err != nil {
		return err
	}
	if err := db.createTaskTagsTable(); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) createLogTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS log;
		CREATE TABLE log(
  			time        VARCHAR(256) NOT NULL,
    		request     VARCHAR(256) NOT NULL,
    		error       text,
    		body        text,
    		query       text,
    		headers     text
			);`)
	return err
}

func (db *DataBase) createUserPrivateTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS user_private CASCADE;
		CREATE TABLE user_private
		(
   	 		user_id  VARCHAR(256) PRIMARY KEY,
    		email    VARCHAR(256) NOT NULL UNIQUE,
    		privacy	 INTEGER NOT NULL,
    		password VARCHAR(256) NOT NULL
		);`)
	return err
}

func (db *DataBase) createUserInfoTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS user_info;
		CREATE TABLE user_info
		(
    		user_id            VARCHAR(256) PRIMARY KEY,
    		first_name         VARCHAR(256) NOT NULL,
    		last_name          VARCHAR(256) NOT NULL,
    		registration_time  VARCHAR(256) NOT NULL,
    		gender             VARCHAR(10)  NOT NULL,
    		online_time        VARCHAR(256) NOT NULL,
    		picture            VARCHAR(512),
    		background_picture VARCHAR(512),
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id) ON DELETE CASCADE
		);`)
	return err
}

func (db *DataBase) createDevelopersTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS developers;
		CREATE TABLE developers
		(
		    user_id VARCHAR(256) NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id) ON DELETE CASCADE
		);`)
	return err
}

func (db *DataBase) createTagsTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS tags;
		CREATE TABLE tags (
  			tag_name VARCHAR (50) PRIMARY KEY,
  			description VARCHAR (256)         
		);`)
	return err
}

func (db *DataBase) createTasksTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS tasks;
		CREATE TABLE tasks (
    		task_name          VARCHAR(50)  PRIMARY KEY,
    		description        VARCHAR(256) NOT NULL,
    		time               VARCHAR(256) NOT NULL,
    		picture            VARCHAR(512),
    		background_picture VARCHAR(512)
		);`)
	return err
}

func (db *DataBase) createUserTagsTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS user_tags;
		CREATE TABLE user_tags (
		    user_id VARCHAR(256) NOT NULL,
    		tag_name VARCHAR(50) NOT NULL,
    		rating INTEGER NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id) ON DELETE CASCADE,
    		FOREIGN KEY (tag_name) REFERENCES tags (tag_name) ON DELETE CASCADE
		);`)
	return err
}

func (db *DataBase) createTaskTagsTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS task_tags;
		CREATE TABLE task_tags (
		    task_name VARCHAR(50) NOT NULL,
    		tag_name VARCHAR(50) NOT NULL,
    		rating INTEGER NOT NULL,
    		FOREIGN KEY (task_name) REFERENCES tasks (task_name) ON DELETE CASCADE,
    		FOREIGN KEY (tag_name) REFERENCES tags (tag_name) ON DELETE CASCADE
		);`)
	return err
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

func (db *DataBase) getUserInfo(result *sql.Rows) (*entities.UserInfo, error) {
	var err error
	ret := &entities.UserInfo{}
	var regTime, onlineTime string

	if err := result.Scan(&ret.UserId, &ret.FirstName, &ret.LastName, &regTime,
		&ret.Gender, &onlineTime, &ret.Picture, &ret.BackgroundPicture); err != nil {
		return nil, err
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

func (db *DataBase) getTaskInfo(result *sql.Rows) (*entities.Task, error) {
	ret := &entities.Task{}
	if err := result.Scan(&ret.Name, &ret.Description, &ret.RecommendedTime,
		&ret.Picture, &ret.BackgroundPicture); err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *DataBase) getTagsList(result *sql.Rows) ([]entities.IdTags, error) {
	ret := make([]entities.IdTags, 0)
	for result.Next() {
		tmp := entities.IdTags{}
		if err := result.Scan(&tmp.Id, &tmp.TagName, &tmp.Rating); err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

func (db *DataBase) errorConstructLong(size, name string) error {
	return errors.New(FieldTooLongError.Error() + size + " bytes for field '" + name + "')")
}

func (db *DataBase) errorConstructValue(err error, name string, values ...string) error {
	var res string
	for _, elem := range values {
		res += elem + "/"
	}
	return errors.New(err.Error() + "'" + name + "' allowed values: " + res)
}

func (db *DataBase) errorConstructNotFound(err error, name string) error {
	return errors.New(err.Error() + "'" + name + "' didn't find")
}
