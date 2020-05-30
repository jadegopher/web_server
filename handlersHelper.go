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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
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
	user, err := handler.DataBase.GetUserInfo(parameters["id"])
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
	if _, in := query[fromField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, fromField)
	}
	if _, in := query[countField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, countField)
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
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

func (handler *Handlers) addTagsToUserHelper(w http.ResponseWriter, r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
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
	if err := handler.validateDeveloperSession(r); err != nil {
		return err
	}
	parameters := mux.Vars(r)
	tags, err := handler.DataBase.GetUserTags(parameters["id"])
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
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

func (handler *Handlers) validateSession(r *http.Request) error {
	sessionId := r.Header.Get(sessionIdField)
	userId := r.Header.Get(userIdField)
	_, err := handler.DataBase.GetUserInfo(userId)
	if err != nil {
		return err
	}
	if sessionId != getSessionId(userId) {
		return invalidSessionError
	}
	return nil
}

func (handler *Handlers) validateDeveloperSession(r *http.Request) error {
	if err := handler.validateSession(r); err != nil {
		return err
	}
	if err := handler.DataBase.CheckDeveloper(r.Header.Get(userIdField)); err != nil {
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

func (handler *Handlers) errorConstructNotFound(err error, name string) error {
	return errors.New(err.Error() + "'" + name + "' didn't find")
}
