package app_error

import "errors"

var (
	InvalidURL       = errors.New("invalid url format")
	SlugAlreadyExist = errors.New("the slug is already in use, please choose another one.")
	InvalidSlug      = errors.New("the slug can only contain letters, numbers, '-' and '_")
)
