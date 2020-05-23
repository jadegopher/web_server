package main

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"web_server/entities"
)

func ToAnswer(element interface{}, err error) *entities.Answer {
	if err != nil {
		args := strings.Split(err.Error(), " ")
		var num = 0
		var tmp = err.Error()
		if len(args) > 0 {
			tmp = strings.Join(args[1:], " ")
			num, err = strconv.Atoi(args[0])
			if err != nil {
				num = 0
				tmp = strings.Join(args, " ")
			}
		}
		return &entities.Answer{
			Error: &entities.Error{
				Code:  num,
				Error: tmp,
			}}
	}
	return &entities.Answer{
		Element: element,
	}
}

func getSessionId(id string) string {
	h := sha256.New()
	h.Write([]byte(id + secret))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
