package store

type Store interface {
	User() UserRepository
	Request() RequestRepository
	Manager() ManagerRepository
}
