package sqlStore

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"database/sql"
)

type RequestRepository struct {
	store *Store
}

func (repo *RequestRepository) FindById(Id int) (*model.Request, error){
	req := &model.Request{}

	if err := repo.store.db.QueryRow("SELECT * FROM public.requests WHERE id = $1",
		Id,
	).Scan(
		&req.Id,
		&req.UserId,
		&req.Message,
		&req.CreateDate,
	); err != nil{

		if err == sql.ErrNoRows{
			return nil, store.ErrRequestNotFound
		}
		return nil, err
	}

	return req, nil
}

func (repo *RequestRepository) NewRequest (newRequest *model.NewRequestRequest) error{
	rez := model.NewRequestResponse{}
	if err := repo.store.db.QueryRow(
		"SELECT addrequest($1, $2)",
		newRequest.UserId,
		newRequest.Message,
	).Scan(&rez.RequestId); err!= nil{
		return err
	}
	return nil
}

func (repo *RequestRepository) CancelRequest (newRequest *model.CancelRequestRequest) (*model.CancelRequestResponse,error){
	rez := &model.CancelRequestResponse{}
	if err := repo.store.db.QueryRow(
		"SELECT cancelrequest($1, $2)",
		newRequest.UserId,
		newRequest.RequestId,
	).Scan(&rez.Request); err!= nil{
		return nil, err
	}
	return rez, nil
}

