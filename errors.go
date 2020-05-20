package main

import "errors"

var (
	invalidSessionError = errors.New("200 session id not valid")
)
