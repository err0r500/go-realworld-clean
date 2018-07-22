package uc_test

import (
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/domain"
	mock "github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProfileGet_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	rick := testData.User("rick")
	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(testData.User("rick").Name).Return(&rick, nil).Times(1)

	retProfile, err := i.GetUCHandler().ProfileGet(testData.User("rick").Name)

	assert.NoError(t, err)
	assert.Equal(t, domain.Profile{User: testData.User("rick")}, *retProfile)

}

func TestProfileGet_fails(t *testing.T) {
	t.Run("err Returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		rick := testData.User("rick")
		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(testData.User("rick").Name).Return(&rick, errors.New("")).Times(1)

		retProfile, err := i.GetUCHandler().ProfileGet(testData.User("rick").Name)

		assert.Error(t, err)
		assert.Nil(t, retProfile)
	})
	t.Run("nil Returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(testData.User("rick").Name).Return(nil, nil).Times(1)

		retProfile, err := i.GetUCHandler().ProfileGet(testData.User("rick").Name)

		assert.Error(t, err)
		assert.Nil(t, retProfile)
	})
	t.Run("wrong User returned", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		jane := testData.User("jane")
		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(testData.User("rick").Name).Return(&jane, nil).Times(1)

		retProfile, err := i.GetUCHandler().ProfileGet(testData.User("rick").Name)

		assert.Error(t, err)
		assert.Nil(t, retProfile)
	})
}
