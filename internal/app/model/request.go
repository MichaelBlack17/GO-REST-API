package model

import "time"

type Request struct {
	Id int
	UserId int
	Message string
	CreateDate time.Time

}
