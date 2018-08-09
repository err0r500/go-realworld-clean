package uc_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	mock "github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEditUser_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedUser := testData.User("rick")
	expectedUser.Email = testData.User("jane").Email
	rick := testData.User("rick")
	jane := testData.User("jane")
	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
	i.UserValidator.EXPECT().CheckUser(expectedUser).Return(nil).Times(1)
	i.UserRW.EXPECT().Save(expectedUser).Return(nil).Times(1)
	i.AuthHandler.EXPECT().GenUserToken(expectedUser.Name).Return("token", nil).Times(1)

	retUser, token, err := i.GetUCHandler().UserEdit(rick.Name, map[domain.UserUpdatableProperty]*string{
		domain.UserEmail: &jane.Email,
		domain.UserBio:   testData.User("jane").Bio, //nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "token", token)
	assert.Equal(t, expectedUser, *retUser)

}
