package model

type ProcessingRequestRequest struct {
	ManagerId 	int 	`json:"manager_id"`
	RequestId 	int 	`json:"request_id"`
}

type ProcessingRequestResponse struct{
	QueueUnit RequestQueue `json:"queue_unit"`
}