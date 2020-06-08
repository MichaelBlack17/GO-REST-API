package sqlStore_test

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store/sqlStore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlStore.TestDB(t, databaseURL)
	defer teardown(`users`)
	s := sqlStore.New(db)
	u := model.User{Name:"Test"}
	err := s.User().Create(&u)
	assert.NoError(t,err)
	assert.NotNil(t,u)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := sqlStore.TestDB(t, databaseURL)
	defer teardown(`users`)
	s := sqlStore.New(db)
	user := model.User{Name:"Test"}
	err := s.User().Create(&user)
	u,err := s.User().FindById(user.Id)
	assert.NoError(t,err)
	assert.NotNil(t,u)
}