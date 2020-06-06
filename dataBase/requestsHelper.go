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
	if err := db.createUserTagsTable(); err != nil {
		return err
	}
	if err := db.createTaskTagsTable(); err != nil {
		return err
	}
	if err := db.createQuestsTable(); err != nil {
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
    		private			   INTEGER NOT NULL,
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
		    user_id VARCHAR(256) PRIMARY KEY,
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id) ON DELETE CASCADE
		);`)
	return err
}

func (db *DataBase) createTagsTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS tags CASCADE;
		CREATE TABLE tags (
  			tag_name VARCHAR (50) PRIMARY KEY,
  			description VARCHAR (256)         
		);`)
	return err
}

func (db *DataBase) createTasksTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS tasks CASCADE;
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

func (db *DataBase) createQuestsTable() error {
	_, err := db.Connection.Exec(`
		DROP TABLE IF EXISTS quests;
		CREATE TABLE quests
		(
    		quest_id      BIGSERIAL PRIMARY KEY,
    		user_id       VARCHAR(256) NOT NULL,
    		task_name     VARCHAR(50),
    		user_opponent VARCHAR(256) NOT NULL,
    		status        INTEGER 	   NOT NULL,
    		start_time    VARCHAR(256),
    		end_time      VARCHAR(256),
    		deadline_time VARCHAR(256),
    		FOREIGN KEY (user_id) REFERENCES user_private (user_id) ON DELETE CASCADE,
    		FOREIGN KEY (user_opponent) REFERENCES user_private (user_id) ON DELETE CASCADE,
    		FOREIGN KEY (task_name) REFERENCES tasks (task_name) ON DELETE CASCADE
		);`)
	return err
}

func (db *DataBase) append(query string, args ...interface{}) error {
	result, err := db.Connection.Exec(query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return AppendError
	}
	return nil
}

func (db *DataBase) getUserInfo(result *sql.Rows) (*entities.UserInfo, error) {
	var err error
	ret := &entities.UserInfo{}
	var regTime, onlineTime string
	var pic, bgPic sql.NullString

	if err := result.Scan(&ret.UserId, &ret.FirstName, &ret.LastName, &regTime,
		&ret.Gender, &onlineTime, &ret.Private, &pic, &bgPic); err != nil {
		return nil, err
	}

	ret.Picture = pic.String
	ret.BackgroundPicture = bgPic.String
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
	var pic, bgPic sql.NullString
	if err := result.Scan(&ret.Name, &ret.Description, &ret.RecommendedTime,
		&pic, &bgPic); err != nil {
		return nil, err
	}
	ret.Picture = pic.String
	ret.BackgroundPicture = bgPic.String
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

func (db *DataBase) getQuestsList(result *sql.Rows) ([]entities.Quest, error) {
	invites := make([]entities.Quest, 0)
	for result.Next() {
		tmp, err := db.getQuest(result)
		if err != nil {
			return nil, err
		}
		invites = append(invites, *tmp)
	}
	return invites, nil
}

func (db *DataBase) getQuest(result *sql.Rows) (*entities.Quest, error) {
	tmp := entities.Quest{}
	var taskName, startTime, endTime, deadlineTime sql.NullString
	if err := result.Scan(&tmp.QuestId, &tmp.UserId, &taskName,
		&tmp.UserOpponent, &tmp.Status, &startTime, &endTime,
		&deadlineTime); err != nil {
		return nil, err
	}
	tmp.TaskName = taskName.String
	tmp.StartTime = startTime.String
	tmp.EndTime = endTime.String
	tmp.DeadlineTime = deadlineTime.String
	return &tmp, nil
}

func (db *DataBase) findQuestWhereOpponent(userId string, questId int, status Status) (*entities.Quest, error) {
	result, err := db.Connection.Query(`SELECT * FROM quests WHERE 
    	quest_id = $1 AND user_opponent = $2 AND status = $3`, questId, userId, status)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, QuestNotFoundError
	}
	tmp, err := db.getQuest(result)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func (db *DataBase) findQuestWhereInitiator(userId string, questId int, status Status) (*entities.Quest, error) {
	result, err := db.Connection.Query(`SELECT * FROM quests WHERE 
    	quest_id = $1 AND user_id = $2 AND status = $3`, questId, userId, status)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, QuestNotFoundError
	}
	tmp, err := db.getQuest(result)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func (db *DataBase) updateQuest(questId int, status Status) error {
	exec, err := db.Connection.Exec(`UPDATE quests SET status = $1 
		WHERE quest_id = $2`, status, questId)
	if err != nil {
		return err
	}
	rows, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return AppendError
	}
	return nil
}

func (db *DataBase) changeToRejected(userId string, questId int) error {
	if _, err := db.findQuestWhereOpponent(userId, questId, Pending); err != nil {
		return err
	}
	if err := db.updateQuest(questId, Rejected); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) changeToNotSelected(userId string, questId int) error {
	quest, err := db.findQuestWhereOpponent(userId, questId, Pending)
	if err != nil {
		return err
	}
	if err := db.updateQuest(questId, NotSelected); err != nil {
		return err
	}
	if err := db.append(`INSERT into quests(user_id, user_opponent, status) 
		VALUES ($1, $2, $3)`, userId, quest.UserId, NotSelected); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) changeToStarted(userId string, questId int) error {
	quest, err := db.findQuestWhereInitiator(userId, questId, Selected)
	if err != nil {
		return err
	}
	if err := db.updateQuest(questId, Started); err != nil {
		return err
	}
	task, err := db.GetTaskInfo(quest.TaskName)
	if err != nil {
		return err
	}
	plus, err := time.Parse(time.RFC1123, task.RecommendedTime)
	if err != nil {
		return err
	}
	startTime := time.Now()
	deadlineTime := startTime.Add(time.Hour * time.Duration(plus.Hour()))
	deadlineTime = deadlineTime.Add(time.Minute * time.Duration(plus.Minute()))
	deadlineTime = deadlineTime.Add(time.Minute * time.Duration(plus.Second()))
	deadlineTime = deadlineTime.Add(time.Minute * time.Duration(plus.Nanosecond()))
	year, month, day := plus.Date()
	deadlineTime = deadlineTime.AddDate(year-1, int(month)-1, day-1)
	exec, err := db.Connection.Exec(`UPDATE quests SET start_time = $1,
        deadline_time = $2, status = $3 WHERE quest_id = $4`,
		startTime.Format(time.RFC1123), deadlineTime.Format(time.RFC1123), Started, questId)
	if err != nil {
		return err
	}
	rows, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return AppendError
	}
	return nil
}

func (db *DataBase) changeToWaiting(userId string, questId int) error {
	quest, err := db.findQuestWhereInitiator(userId, questId, Started)
	if err != nil {
		return err
	}
	timeNow, err := db.checkExpired(quest)
	if err != nil {
		return err
	}
	exec, err := db.Connection.Exec(`UPDATE quests SET end_time = $1, status = $2 WHERE quest_id = $3`,
		timeNow.Format(time.RFC1123), Waiting, questId)
	if err != nil {
		return err
	}
	rows, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return AppendError
	}
	return nil
}

func (db *DataBase) changeToAccepted(userId string, questId int) error {
	_, err := db.findQuestWhereOpponent(userId, questId, Waiting)
	if err != nil {
		return err
	}
	if err := db.updateQuest(questId, Accepted); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) changeToNotAccepted(userId string, questId int) error {
	_, err := db.findQuestWhereOpponent(userId, questId, Waiting)
	if err != nil {
		return err
	}
	if err := db.updateQuest(questId, NotAccepted); err != nil {
		return err
	}
	return nil
}

func (db *DataBase) checkExpired(quest *entities.Quest) (*time.Time, error) {
	if quest.Status == int(Expired) {
		return nil, ExpiredTaskError
	}
	timeNow := time.Now()
	deadline, err := time.Parse(time.RFC1123, quest.DeadlineTime)
	if err != nil {
		return nil, err
	}
	if timeNow.After(deadline) {
		if err := db.updateQuest(quest.QuestId, Expired); err != nil {
			return nil, err
		}
		return nil, ExpiredTaskError
	}
	return &timeNow, nil
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
