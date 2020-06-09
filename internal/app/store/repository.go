package store

import "GO-REST-API/internal/app/model"

type UserRepository interface {
	Create(user *model.User) error
	FindById(Id int) (*model.User, error)
}

type RequestRepository interface {
	NewRequest (newRequest *model.NewRequestRequest) error
}
