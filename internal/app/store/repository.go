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
}
