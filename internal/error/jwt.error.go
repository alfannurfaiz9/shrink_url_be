package app_error

import "errors"

var ErrMissingKey = errors.New("jwt key not found")
