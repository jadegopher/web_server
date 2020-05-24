package dataBase

import (
	"database/sql"
	"errors"
	"time"
	"web_server/entities"
)

func (db *DataBase) createNewDataBase() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS log;
		DROP TABLE IF EXISTS user_info;
		DROP TABLE IF EXISTS developers;
		DROP TABLE IF EXISTS user_private;
		DROP TABLE IF EXISTS tags;
		CREATE TABLE log(
  			time        VARCHAR(256) NOT NULL,
    		request     VARCHAR(256) NOT NULL,
    		error       text,
    		body        text,
    		query       text,
    		headers     text
			);
		CREATE TABLE user_private
		(
   	 		user_id  VARCHAR(256) NOT NULL UNIQUE PRIMARY KEY,
    		email    VARCHAR(256) NOT NULL UNIQUE,
    		password VARCHAR(256) NOT NULL
		);
		CREATE TABLE developers
		(
		    user_id VARCHAR(256) NOT NULL,
		    PRIMARY KEY (user_id),
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id)
		);
		CREATE TABLE user_info
		(
    		user_id            VARCHAR(256) NOT NULL UNIQUE,
    		first_name         VARCHAR(256) NOT NULL,
    		last_name          VARCHAR(256) NOT NULL,
    		registration_time  VARCHAR(256) NOT NULL,
    		gender             VARCHAR(10)  NOT NULL,
    		online_time        VARCHAR(256) NOT NULL,
    		picture            VARCHAR(512),
    		background_picture VARCHAR(512),
    		PRIMARY KEY (user_id),
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id)
		);
		CREATE TABLE tags (
  			tag_name VARCHAR (50) NOT NULL UNIQUE,
  			description VARCHAR (256)         
		);`)
	if err != nil {
		return err
	}
	return nil
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
