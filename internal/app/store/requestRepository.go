package store

import "GO-REST-API/internal/app/model"

type RequestRepository struct {
	store *Store
}

func (repo *RequestRepository) Create (newRequest *model.NewRequestRequest)(*model.NewRequestResponse, error){
	rez := model.NewRequestResponse{}
	if err := repo.store.db.QueryRow(
		"SELECT addrequest($1)",
		newRequest.Message,
	).Scan(&rez.RequestId); err!= nil{
		return nil, err
	}
	return &rez, nil
}

