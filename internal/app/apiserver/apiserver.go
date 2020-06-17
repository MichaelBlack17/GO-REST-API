package apiserver

import (
	"GO-REST-API/internal/app/store"
	"GO-REST-API/internal/app/store/sqlStore"
	"database/sql"
	"net/http"
	"strconv"
)

func Start(config *Config) error{
	db, err := newDB(config.DatabaseURL)

	if err != nil{
		return err
	}

	defer db.Close()
	command := "set glb.queue_length to " + strconv.Itoa(config.QueueLength)
	if _,err := db.Exec(command); err != nil{
		return  store.ErrParamNotFound
	}

	command = "set glb.valid_time to " + strconv.Itoa(config.ValidTimeOut)
	if _,err := db.Exec(command); err != nil{
		return  store.ErrParamNotFound
	}

	store := sqlStore.New(db)
	svr := newServer(store)

	go store.Request().StartQueryManagement(config.ValidTimeOut)

	return http.ListenAndServe(config.BindAddr, svr)
}

func newDB(databaseURL string)(*sql.DB, error){
	db, err := sql.Open("postgres",databaseURL)
	if err != nil{
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return db, nil
}