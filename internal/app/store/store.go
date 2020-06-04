package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	)

type Store struct {
	db *sql.DB
	userRepository *UserRepository
	requestRepository *RequestRepository

}

func New(db *sql.DB) *Store {
	return &Store{
		db:db,
	}
}


func (s *Store) User() *UserRepository {
	if s.userRepository != nil{
		return s.userRepository
	}
	s.userRepository = &UserRepository{store: s}
	return s.userRepository
}

func (s *Store) Request() *RequestRepository {
	if s.requestRepository != nil{
		return s.requestRepository
	}
	s.requestRepository = &RequestRepository{store: s}
	return s.requestRepository
}