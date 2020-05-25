package main

import (
	"net/http"
	"web_server/dataBase"
)

type Handlers struct {
	DataBase *dataBase.DataBase
	Log      bool
}

func NewHandlers(db *dataBase.DataBase, log bool) *Handlers {
	return &Handlers{DataBase: db, Log: log}
}

func (handler *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.registrationHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "Registration", err)
	}
}

func (handler *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.loginHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "Login", err)
	}
}

func (handler *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getUserInfoHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetUserInfo", err)
	}
}

func (handler *Handlers) SearchUser(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.searchUserHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "SearchUser", err)
	}
}

func (handler *Handlers) GetDeveloperAccount(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getDeveloperAccountHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetDeveloperAccount", err)
	}
}

func (handler *Handlers) PostTag(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.postTagHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "PostTag", err)
	}
}

func (handler *Handlers) PostTask(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.postTaskHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "PostTask", err)
	}
}
