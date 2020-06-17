package store

import "GO-REST-API/internal/app/model"

type UserRepository interface {
	Create(user *model.User) error
	FindById(Id int) (*model.User, error)
}

type RequestRepository interface {
	NewRequest (newRequest *model.NewRequestRequest) error
	CancelRequest (Request *model.CancelRequestRequest) (*model.CancelRequestResponse, error)
	FindById(Id int) (*model.Request, error)
	FindByUserAndReqId(UserId int, ReqId int) (*model.Request, error)
	AllUserRequests (req *model.AllUserRequestsRequest) (*model.AllUserRequestsResponse,error)
	ProcessingRequest (req *model.ProcessingRequestRequest) (*model.ProcessingRequestResponse,error)
	CancelProcessingRequest (req *model.CancelProcessingRequestRequest) (*model.CancelProcessingRequestResponse,error)
	AllManagerRequests (req *model.AllManagerRequestsRequest) (*model.AllManagerRequestsResponse,error)
	StartQueryManagement (mins int) error
}

type ManagerRepository interface {
	FindById(Id int) (*model.Manager, error)
	FindByManagerAndReqId(ManagerId int, ReqId int) (*model.RequestQueue, error)
}
