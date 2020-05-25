package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"web_server/dataBase"
	"web_server/entities"
)

func getConfig(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := make(map[string]interface{})
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	tmp := &dataBase.DBConfig{}
	ret[dataBaseField], err = getModule(cfg[dataBaseField], tmp)
	ret[logField] = cfg[logField]
	return ret, nil
}

func getModule(in, out interface{}) (interface{}, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func toAnswer(element interface{}, err error) *entities.Answer {
	if err != nil {
		tmp, num := parseError(err)
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

func parseError(err error) (string, int) {
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
	return tmp, num
}

func getSessionId(id string) string {
	h := sha256.New()
	h.Write([]byte(id + secretSessionId))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
