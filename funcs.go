package main

import (
	"crypto/sha256"
	"encoding/base64"
	"web_server/entities"
)

func ToAnswer(element interface{}, err error) *entities.Answer {
	var tmp string
	if err != nil {
		tmp = err.Error()
	} else {
		tmp = "nil"
	}
	return &entities.Answer{
		Element: element,
		Error:   tmp,
	}
}

func getSessionId(id string) string {
	h := sha256.New()
	h.Write([]byte(id + "secretPassword"))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
