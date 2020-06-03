package store_test

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := store.TestDB(t, databaseURL)
	defer teardown(`users`)
	s := store.New(db)
	u,err := s.User().Create(&model.User{Name:"Test"})
	assert.NoError(t,err)
	assert.NotNil(t,u)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := store.TestDB(t, databaseURL)
	defer teardown(`users`)
	s := store.New(db)
	user,err := s.User().Create(&model.User{Name:"Test"})
	u,err := s.User().FindById(user.Id)
	assert.NoError(t,err)
	assert.NotNil(t,u)
}