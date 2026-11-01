package app_error

import "errors"

var (
	ErrConnectDatabase = errors.New("cannot connect to database")
	EmailAlreadyExist  = errors.New("email already exist")
)
