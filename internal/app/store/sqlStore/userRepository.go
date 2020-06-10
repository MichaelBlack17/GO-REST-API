package sqlStore

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (repo *UserRepository) Create (user *model.User) error {
	if err := repo.store.db.QueryRow(
		"INSERT INTO public.users(Name) VALUES ($1) RETURNING Id",
		user.Name,
		).Scan(&user.Id); err!= nil{
		return err
	}
	return  nil
}

func (repo *UserRepository) FindById(Id int) (*model.User, error){
	user := &model.User{}

	if err := repo.store.db.QueryRow("SELECT id, Name FROM public.users WHERE id = $1",
		Id,
	).Scan(
		&user.Id,
		&user.Name,
	); err != nil{

		if err == sql.ErrNoRows{
			return nil, store.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}