package sqlStore

import "GO-REST-API/internal/app/model"

type RequestRepository struct {
	store *Store
}

func (repo *RequestRepository) NewRequest (newRequest *model.NewRequestRequest) error{
	rez := model.NewRequestResponse{}
	if err := repo.store.db.QueryRow(
		"SELECT addrequest($1)",
		newRequest.Message,
	).Scan(&rez.RequestId); err!= nil{
		return err
	}
	return nil
}

