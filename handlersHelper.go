package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"web_server/entities"
)

func (handler *Handlers) registrationHelper(w http.ResponseWriter, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	userInfo := &entities.Registration{}
	if err = json.Unmarshal(data, userInfo); err != nil {
		return err
	}
	userInfo.UserInfo.RegistrationTime = time.Now()
	if err = handler.DataBase.AddUser(userInfo); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) loginHelper(w http.ResponseWriter, r *http.Request) error {
	data, err := handler.copyBody(r)
	if err != nil {
		return err
	}
	userInfo := &entities.Login{}
	if err = json.Unmarshal(data, userInfo); err != nil {
		return err
	}
	userId, err := handler.DataBase.Login(userInfo)
	if err != nil {
		return err
	}
	w.Header().Set(sessionIdField, getSessionId(userId))
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getUserInfoHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	//fmt.Println(parameters["id"])
	var self bool
	if parameters["id"] == r.Header.Get(userIdField) {
		self = true
	} else {
		self = false
	}
	user, err := handler.DataBase.GetUserInfo(parameters["id"], self)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(&user, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) searchUserHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	query := r.URL.Query()
	if _, in := query[queryField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, queryField)
	}
	if err := handler.validateFromCountFields(query); err != nil {
		return err
	}
	users, err := handler.DataBase.SearchUser(r.Header.Get(userIdField), query[queryField][0], query[fromField][0], query[countField][0])
	if err != nil {
		return err
	}
	if err = json.NewEncoder(w).Encode(toAnswer(users, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) deleteUserHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	data, err := handler.copyBody(r)
	if err != nil {
		return err
	}
	userInfo := &entities.Login{}
	if err = json.Unmarshal(data, userInfo); err != nil {
		return err
	}
	userId, err := handler.DataBase.Login(userInfo)
	if err != nil {
		return err
	}
	if userId != r.Header.Get(userIdField) {
		return credError
	}
	if err := handler.DataBase.DeleteUser(userId); err != nil {
		return err
	}
	if err = json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getTagsHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	query := r.URL.Query()
	if err := handler.validateFromCountFields(query); err != nil {
		return err
	}
	tags, err := handler.DataBase.GetTags(query[fromField][0], query[countField][0])
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(tags, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) addTagsToUserHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	data, err := handler.copyBody(r)
	if err != nil {
		return err
	}
	var tags []entities.Tag
	if err = json.Unmarshal(data, &tags); err != nil {
		return err
	}
	if err = handler.DataBase.AddTagsToUser(r.Header.Get(userIdField), tags); err != nil {
		return err
	}
	if err = json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getUserTagsHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	var self bool
	if parameters["id"] == r.Header.Get(userIdField) {
		self = true
	} else {
		self = false
	}
	tags, err := handler.DataBase.GetUserTags(parameters["id"], self)
	if err != nil {
		return err
	}
	if err = json.NewEncoder(w).Encode(toAnswer(tags, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getTaskInfoHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	task, err := handler.DataBase.GetTaskInfo(parameters["taskName"])
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(task, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getTaskTagsHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	tags, err := handler.DataBase.GetTaskTags(parameters["taskName"])
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(tags, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) inviteUserHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	if err := handler.DataBase.InviteUser(r.Header.Get(userIdField), parameters["id"]); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getInvitesHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	invites, err := handler.DataBase.GetInvites(r.Header.Get(userIdField))
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(invites, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getQuestsHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	query := r.URL.Query()
	if _, in := query[statusField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, statusField)
	}
	invites, err := handler.DataBase.GetQuests(r.Header.Get(userIdField), query[statusField][0])
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(invites, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) getDeveloperAccountHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	secret := r.Header.Get(developerField)
	if secret != secretDeveloper {
		return developerSecretError
	}
	if err := handler.DataBase.AddDeveloper(r.Header.Get(userIdField)); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) postTagHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateDeveloperSession(r); err != nil {
		return err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	tagInfo := &entities.Tag{}
	if err = json.Unmarshal(data, tagInfo); err != nil {
		return err
	}
	if err = handler.DataBase.AddTag(tagInfo); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) postTaskHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateDeveloperSession(r); err != nil {
		return err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	taskInfo := &entities.Task{}
	if err = json.Unmarshal(data, taskInfo); err != nil {
		return err
	}
	if _, err := time.Parse(taskInfo.RecommendedTime, time.Time{}.Format(time.RFC1123)); err != nil {
		return err
	}
	if err = handler.DataBase.AddTask(taskInfo); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) addTagsToTaskHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateDeveloperSession(r); err != nil {
		return err
	}
	data, err := handler.copyBody(r)
	if err != nil {
		return err
	}
	var tags []entities.Tag
	if err := json.Unmarshal(data, &tags); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	if err := handler.DataBase.AddTagsToTask(parameters["taskName"], tags); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(toAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) defaultErrorResponse(w http.ResponseWriter, err error) {
	data, err := json.Marshal(toAnswer(nil, err))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}
	http.Error(w, string(data), http.StatusBadRequest)
}

func (handler *Handlers) addLog(r *http.Request, reqName string, userErr error) {
	headers, err := json.Marshal(r.Header)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	query, err := json.Marshal(r.URL.Query())
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	tmp := &entities.Log{
		Time:    time.Now().Format(time.RFC1123),
		Request: reqName,
		Error:   "",
		Body:    string(body),
		Query:   string(query),
		Headers: string(headers),
	}
	if userErr != nil {
		tmp.Error = userErr.Error()
	}
	if err = handler.DataBase.LogAdd(tmp); err != nil {
		log.Fatal(err.Error())
		return
	}
}

func (handler *Handlers) copyBody(r *http.Request) ([]byte, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = r.Body.Close()
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return data, nil
}

func (handler *Handlers) errorConstructNotFound(err error, name string) error {
	return errors.New(err.Error() + "'" + name + "' didn't find")
}
