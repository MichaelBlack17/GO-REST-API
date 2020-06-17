package model

import "time"

type RequestQueue struct {
	Id 			int 		`json:"id"`
	ManagerId 	int 		`json:"manager_id"`
	RequestId 	int 		`json:"request_id"`
	Status 		int 		`json:"status"`
	ValidTime 	time.Time	`json:"valid_time"`
}
