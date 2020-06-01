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

func (handler *Handlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.deleteUserHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "DeleteUser", err)
	}
}

func (handler *Handlers) GetTags(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getTagsHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetTags", err)
	}
}

func (handler *Handlers) AddTagsToUser(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.addTagsToUserHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "AddTagsToUser", err)
	}
}

func (handler *Handlers) GetUserTags(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getUserTagsHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetUserTags", err)
	}
}

func (handler *Handlers) GetTaskInfo(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getTaskInfoHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetTaskInfo", err)
	}
}

func (handler *Handlers) GetTaskTags(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getTaskTagsHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetTaskTags", err)
	}
}

func (handler *Handlers) InviteUser(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.inviteUserHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "InviteUser", err)
	}
}

func (handler *Handlers) GetInvites(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getInvitesHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetInvites", err)
	}
}

func (handler *Handlers) GetQuests(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.getQuestsHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "GetQuests", err)
	}
}

func (handler *Handlers) ChangeQuestStatus(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.changeQuestStatusHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "ChangeQuestStatus", err)
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

func (handler *Handlers) AddTagsToTask(w http.ResponseWriter, r *http.Request) {
	var err error
	if err = handler.addTagsToTaskHelper(w, r); err != nil {
		handler.defaultErrorResponse(w, err)
	}
	if handler.Log {
		handler.addLog(r, "AddTagsToTask", err)
	}
}
