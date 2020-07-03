package apiserver

import (
	"GO-REST-API/internal/app/store/sqlStore"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"context"
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
		if err := InitGlobalFunctions(db, config.QueueLength, config.ValidTimeOut); err != nil {
			svr.logger.Errorf("Error: %v", err)
		}
		for {

			<-time.After(time.Minute * time.Duration(config.ValidTimeOut))

			if err := store.Request().StartQueryManagement(); err != nil {
				svr.logger.Errorf("Error: %v", err)
			}
			svr.logger.Info("queue management finished")

		}
	}()

	fmt.Println("server starting...")
	//return http.ListenAndServe(config.BindAddr, svr)
	server := http.Server{
		Addr: config.BindAddr,
		Handler: svr.router,
	}

	server.ListenAndServe()
	// ...
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	return nil
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
		return fmt.Errorf("error in set queue_length database parametr: %v", err)
	}

	command = "set glb.valid_time to " + strconv.Itoa(ValidTimeOut)
	if _, err := db.Exec(command); err != nil {
		return fmt.Errorf("error in set valid_time database parametr: %v", err)
	}
	return nil
}
