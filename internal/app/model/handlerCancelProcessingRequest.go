package model

type CancelProcessingRequestRequest struct {
	ManagerId int `json:"manager_id"`
	RequestId int `json:"request_id"`
}

type CancelProcessingRequestResponse struct {
	QueueUnit RequestQueue `json:"queue_unit"`
}
