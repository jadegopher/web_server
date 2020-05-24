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
	"web_server/dataBase"
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
	if err = handler.validateRegFields(userInfo); err != nil {
		return err
	}
	userInfo.RegistrationTime = time.Now()
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
	userInfo := &entities.UserPrivate{}
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
		return handler.errorConstructField(fieldNotFoundError, queryField)
	}
	if _, in := query[fromField]; !in {
		return handler.errorConstructField(fieldNotFoundError, fromField)
	}
	if _, in := query[countField]; !in {
		return handler.errorConstructField(fieldNotFoundError, countField)
	}
	users, err := handler.DataBase.SearchUser(query[queryField][0], query[fromField][0], query[countField][0])
	if err != nil {
		return err
	}
	if err = json.NewEncoder(w).Encode(toAnswer(users, nil)); err != nil {
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

func (handler *Handlers) defaultErrorResponse(w http.ResponseWriter, err error) {
	data, err := json.Marshal(toAnswer(nil, err))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}
	http.Error(w, string(data), http.StatusBadRequest)
}

func (handler *Handlers) validateRegFields(userInfo *entities.Registration) error {
	if userInfo.UserId == "" {
		return handler.errorConstructField(fieldNotFoundError, userIdField)
	}
	if userInfo.Email == "" {
		return handler.errorConstructField(fieldNotFoundError, emailField)
	}
	if userInfo.FirstName == "" {
		return handler.errorConstructField(fieldNotFoundError, firstNameField)
	}
	if userInfo.LastName == "" {
		return handler.errorConstructField(fieldNotFoundError, lastNameField)
	}
	if userInfo.Password == "" {
		return handler.errorConstructField(fieldNotFoundError, passwordField)
	}
	return nil
}

func (handler *Handlers) errorConstructField(err error, add string) error {
	return errors.New(err.Error() + "'" + add + "' " + "doesn't exist")
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
	tmp := &dataBase.Log{
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
