package uc_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserGet_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	authToken := "anyToken"
	rick := testData.User("rick")
	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
	i.AuthHandler.EXPECT().GenUserToken(rick.Name).Return(authToken, nil).Times(1)

	retUser, token, err := i.GetUCHandler().UserGet(rick.Name)
	assert.NoError(t, err)
	assert.Equal(t, rick, *retUser)
	assert.Equal(t, authToken, token)
}
