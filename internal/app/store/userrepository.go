package store

import "GO-REST-API/internal/app/model"

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) Create (user *model.User)(*model.User, error){
	if err := repo.store.db.QueryRow(
		"INSERT INTO Users(Name) VALUES ($1) RETURNING Id",
		user.Name,
		).Scan(&user.Id); err!= nil{
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindById(Id int64)	(*model.User, error){
	return nil,nil
}