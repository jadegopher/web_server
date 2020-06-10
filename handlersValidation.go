package main

import (
	"net/http"
	"net/url"
)

func (handler *Handlers) validateSession(r *http.Request) error {
	sessionId := r.Header.Get(sessionIdField)
	userId := r.Header.Get(userIdField)
	_, err := handler.DataBase.GetUserInfo(userId, true)
	if err != nil {
		return err
	}
	if sessionId != getSecret(userId, secretSessionId) {
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

func (handler *Handlers) validateFromCountFields(query url.Values) error {
	if _, in := query[fromField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, fromField)
	}
	if _, in := query[countField]; !in {
		return handler.errorConstructNotFound(fieldNotFoundError, countField)
	}
	return nil
}
