package model

import "time"

type RequestQueue struct {
	Id int
	ManagerId int
	RequestId int
	Status int
	ValidTime time.Time
}
