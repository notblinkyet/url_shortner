package errors

import "errors"

var (
	ErrAliasAlreadyUse error = errors.New("this alias is already use, generate another")
	ErrFindPass        error = errors.New("can't find POSTGRES_PASS in env")
	ErrAlreadyExist    error = errors.New("url is already exists")
	ErrAliaceDontUse   error = errors.New("can't find aliace")
)
