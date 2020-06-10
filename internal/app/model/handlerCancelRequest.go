package model

type CancelRequestRequest struct {
	UserId 	int 	`json:"user_id"`
	RequestId int 	`json:"request_id"`
}

type CancelRequestResponse struct{
	Request Request `json:"request"`
}
