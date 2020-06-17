package model

type AllManagerRequestsRequest struct {
	ManagerId 	int 	`json:"manager_id"`
}

type AllManagerRequestsResponse struct {
	RequestList []Request `json:"request_list"`
}