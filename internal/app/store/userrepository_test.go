package store_test

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	var s, teardown = store.TestStore(t, databaseURL)
	defer teardown(`public."Users"`)

	u,err := s.User().Create(&model.User{Name:"Michael",})
	assert.NoError(t,err)
	assert.NotNil(t,u)
}