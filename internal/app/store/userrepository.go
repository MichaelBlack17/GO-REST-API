package store

import "GO-REST-API/internal/app/model"

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) Create (user *model.User)(*model.User, error){
	if err := repo.store.db.QueryRow(
		"INSERT INTO public.users(Name) VALUES ($1) RETURNING Id",
		user.Name,
		).Scan(&user.Id); err!= nil{
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindById(Id int) (*model.User, error){
	user := &model.User{}

	if err := repo.store.db.QueryRow("SELECT Id, Name FROM public.users WHERE Id = $1",
		Id,
	).Scan(&user.Id, &user.Name); err != nil{
		return nil, err
	}

	return user, nil
}