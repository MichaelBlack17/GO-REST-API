package sqlStore

import (
	"GO-REST-API/internal/app/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	requestRepository *RequestRepository
	managerRepository *ManagerRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{store: s}
	return s.userRepository
}

func (s *Store) Request() store.RequestRepository {
	if s.requestRepository != nil {
		return s.requestRepository
	}
	s.requestRepository = &RequestRepository{store: s}
	return s.requestRepository
}

func (s *Store) Manager() store.ManagerRepository {
	if s.managerRepository != nil {
		return s.managerRepository
	}
	s.managerRepository = &ManagerRepository{store: s}
	return s.managerRepository
}
