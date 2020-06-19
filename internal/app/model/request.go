package model

import "time"

type Request struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Message    string    `json:"message"`
	CreateDate time.Time `json:"create_date"`
}
