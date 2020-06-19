package model

type NewRequestRequest struct {
	UserId  int    `json:"user_id"`
	Message string `json:"message"`
}

type NewRequestResponse struct {
	RequestId int `json:"request_id"`
}
