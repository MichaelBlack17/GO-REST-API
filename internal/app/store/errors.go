package store

import "errors"

var (
	ErrRecordNotFound 	= errors.New("record not found")
	ErrUserNotFound 	= errors.New("user not found")
	ErrRequestNotFound 	= errors.New("request not found")
	ErrParamNotFound 	= errors.New("DB param not found")
	ErrQueueManage 	= errors.New("error in management queue")
)
