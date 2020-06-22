package apiserver

import (
	"GO-REST-API/internal/app/store/sqlStore"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)

	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlStore.New(db)
	svr := newServer(store)

	go func() {
		for {
			select {

			case <-time.After(time.Minute * time.Duration(config.ValidTimeOut)):
				if err := InitGlobalFunctions(db, config.QueueLength, config.ValidTimeOut); err != nil {
					fmt.Errorf("error", err)
				}

				if err := store.Request().StartQueryManagement(); err != nil {
					fmt.Errorf("error", err)
				}
				svr.logger.Info("queue management finished")

			}

		}
	}()

	fmt.Println("server starting...")
	return http.ListenAndServe(config.BindAddr, svr)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitGlobalFunctions(db *sql.DB, QueueLength int, ValidTimeOut int) error {
	command := "set glb.queue_length to " + strconv.Itoa(QueueLength)
	if _, err := db.Exec(command); err != nil {
		return fmt.Errorf("error in set queue_length database parametr")
	}

	command = "set glb.valid_time to " + strconv.Itoa(ValidTimeOut)
	if _, err := db.Exec(command); err != nil {
		return fmt.Errorf("error in set valid_time database parametr")
	}
	return nil
}
