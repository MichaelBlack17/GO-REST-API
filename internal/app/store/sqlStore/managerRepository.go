package sqlStore

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"database/sql"
)

type ManagerRepository struct {
	store *Store
}

func (repo *ManagerRepository) FindById(Id int) (*model.Manager, error) {
	req := &model.Manager{}

	if err := repo.store.db.QueryRow("SELECT * FROM public.managers WHERE id = $1",
		Id,
	).Scan(
		&req.Id,
		&req.Name,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, store.ErrRequestNotFound
		}
		return nil, err
	}
	return req, nil
}

func (repo *ManagerRepository) FindByManagerAndReqId(ManagerId int, ReqId int) (*model.RequestQueue, error) {
	req := &model.RequestQueue{}

	if err := repo.store.db.QueryRow("SELECT * FROM public.requestqueue WHERE request_id = $1 AND manager_id = $2 AND status <> 2",
		ReqId,
		ManagerId,
	).Scan(
		&req.Id,
		&req.ManagerId,
		&req.RequestId,
		&req.Status,
		&req.ValidTime,
	); err != nil {

		if err == sql.ErrNoRows {
			return nil, store.ErrRequestNotFound
		}
		return nil, err
	}
	return req, nil
}
