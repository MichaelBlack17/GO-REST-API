package apiserver

import (
	"GO-REST-API/internal/app/store/sqlStore"
	"database/sql"
	"net/http"
)

func Start(config *Config) error{
	db, err := newDB(config.DatabaseURL)

	if err != nil{
		return err
	}

	defer db.Close()

	store := sqlStore.New(db)
	svr := newServer(store)
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