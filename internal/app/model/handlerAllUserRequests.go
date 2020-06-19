package model

type AllUserRequestsRequest struct {
	UserId int `json:"user_id"`
}

type AllUserRequestsResponse struct {
	RequestList []Request `json:"request_list"`
}
