package main

import (
	"net/http"
	"web_server/dataBase"
)

type Handlers struct {
	DataBase *dataBase.DataBase
}

func NewHandlers(db *dataBase.DataBase) *Handlers {
	return &Handlers{DataBase: db}
}

func (handler *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	if err := handler.registrationHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
}

func (handler *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if err := handler.loginHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
}

func (handler *Handlers) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if err := handler.getUserInfoHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
}
