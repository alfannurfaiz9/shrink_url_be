package app_error

import "errors"

var (
	IncorrectHashFormat  = errors.New("incorect hash format")
	IncorrectHashType    = errors.New("incorect hash type")
	IncorrectHashVersion = errors.New("incorect hash version")
	MismatchHash         = errors.New("mismatch hash")
)
