package store

import "errors"

var (
	ErrRecordNotFound 	= errors.New("record not found")
	ErrUserNotFound 	= errors.New("user not found")
	ErrRequestNotFound 	= errors.New("request not found")
)
