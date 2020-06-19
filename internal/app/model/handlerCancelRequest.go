package model

type CancelRequestRequest struct {
	UserId    int `json:"user_id"`
	RequestId int `json:"request_id"`
}

type CancelRequestResponse struct {
	QueueRow RequestQueue `json:"queue_row"`
}
