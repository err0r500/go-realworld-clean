package userRW_test

import (
	"sync"
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	userRW "github.com/err0r500/go-realworld-clean/implem/memory.userRW"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/stretchr/testify/assert"
)

func TestRw_Create(t *testing.T) {
	rw := userRW.New()

	happyInsert := func(t *testing.T, toInsert domain.User) {
		returnedUser, err := rw.Create(toInsert.Name, toInsert.Email, toInsert.Password)
		assert.NoError(t, err)
		assert.Equal(t, toInsert.Name, returnedUser.Name)
		assert.Equal(t, toInsert.Email, returnedUser.Email)
		assert.Equal(t, toInsert.Password, returnedUser.Password)
	}

	faillingInsert := func(t *testing.T, toInsert domain.User) {
		returnedUser, err := rw.Create(toInsert.Name, toInsert.Email, toInsert.Password)
		assert.Error(t, err)
		assert.Nil(t, returnedUser)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func(t *testing.T, toInsert domain.User) {
		defer wg.Done()
		happyInsert(t, toInsert)
		faillingInsert(t, toInsert)
	}(t, testData.User("jane"))
	go func(t *testing.T, toInsert domain.User) {
		defer wg.Done()
		happyInsert(t, toInsert)
		faillingInsert(t, toInsert)
	}(t, testData.User("rick"))
	wg.Wait()
}

func TestRw_GetByEmailAndPassword(t *testing.T) {
	rw := userRW.New()
	jane := testData.User("jane")

	rw.Create(jane.Name, jane.Email, jane.Password)
	rw.Save(jane)

	user, err := rw.GetByEmailAndPassword(jane.Email, jane.Password)

	assert.NoError(t, err)
	assert.Equal(t, jane, *user)
}
