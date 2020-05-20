package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
	"web_server/dataBase"
	"web_server/entities"
)

type Handlers struct {
	DataBase *dataBase.DataBase
}

func NewHandlers(db *dataBase.DataBase) *Handlers {
	return &Handlers{DataBase: db}
}

func (handler *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	if err := handler.registrationHelper(w, r); err != nil {
		json.NewEncoder(w).Encode(ToAnswer(fail, err))
	}
}

func (handler *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if err := handler.loginHelper(w, r); err != nil {
		json.NewEncoder(w).Encode(ToAnswer(fail, err))
	}
}

func (handler *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if err := handler.getUserInfoHelper(w, r); err != nil {
		json.NewEncoder(w).Encode(ToAnswer(fail, err))
	}
}

func (handler *Handlers) registrationHelper(w http.ResponseWriter, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	userInfo := &entities.Registration{}
	if err = json.Unmarshal(data, userInfo); err != nil {
		return err
	}
	userInfo.RegistrationTime = time.Now()
	if err = handler.DataBase.AddUser(userInfo); err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(ToAnswer(success, nil)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) loginHelper(w http.ResponseWriter, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	userInfo := &entities.UserPrivate{}
	if err = json.Unmarshal(data, userInfo); err != nil {
		return err
	}
	if err = handler.DataBase.Login(userInfo); err != nil {
		return err
	}
	w.Header().Set("sessionId", getSessionId(userInfo.UserId))
	if err := json.NewEncoder(w).Encode(ToAnswer(success, err)); err != nil {
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
	if err := json.NewEncoder(w).Encode(ToAnswer(&user, err)); err != nil {
		return err
	}
	return nil
}

func (handler *Handlers) validateSession(r *http.Request) error {
	sessionId := r.Header.Get("sessionId")
	userId := r.Header.Get("userId")
	_, err := handler.DataBase.GetUserInfo(userId)
	if err != nil {
		return err
	}
	if sessionId != getSessionId(userId) {
		return invalidSessionError
	}
	return nil
}
