package sqlStore

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"database/sql"
	"encoding/json"
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

func (repo *RequestRepository) FindByUserAndReqId(UserId int, ReqId int) (*model.Request, error){
	req := &model.Request{}

	if err := repo.store.db.QueryRow("SELECT * FROM public.requests WHERE id = $1 AND user_id = $2",
		ReqId,
		UserId,
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
	b := []byte(`{}`)
	if err := repo.store.db.QueryRow(
		"SELECT cancelrequest($1, $2)",
		newRequest.UserId,
		newRequest.RequestId,
	).Scan(
		&b,
		); err!= nil{
		return nil, err
	}

	if err := json.Unmarshal(b, &rez.QueueRow); err!= nil{
		return nil, err
	}

	return rez, nil
}

func (repo *RequestRepository) AllUserRequests (req *model.AllUserRequestsRequest) (*model.AllUserRequestsResponse,error){
	rez := &model.AllUserRequestsResponse{}
	b := []byte(`{}`)
	if err := repo.store.db.QueryRow(
		"SELECT json_agg(r.*) as tags FROM (SELECT * FROM public.requests WHERE user_id = $1) as r",
		req.UserId,
	).Scan(
		&b,
	); err!= nil{
		return nil, err
	}
	
	if err := json.Unmarshal(b, &rez.RequestList); err!= nil{
		return nil, err
	}

	return rez, nil
}

