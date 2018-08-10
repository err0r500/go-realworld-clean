package uc_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProfileUpdateFollow_happyCase(t *testing.T) {
	t.Run("add followee", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		toSave := testData.User("rick")
		toSave.FollowIDs = []string{testData.User("jane").Name}
		rick := testData.User("rick")

		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
		i.UserRW.EXPECT().Save(toSave).Return(nil).Times(1)

		retProfile, err := i.GetUCHandler().ProfileUpdateFollow(rick.Name, testData.User("jane").Name, true)
		assert.NoError(t, err)
		assert.Equal(t, toSave, *retProfile)
	})

	t.Run("remove followee", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		withFollowee := testData.User("rick")
		withFollowee.FollowIDs = []string{testData.User("jane").Name}

		i := mock.NewMockedInteractor(mockCtrl)
		i.UserRW.EXPECT().GetByName(testData.User("rick").Name).Return(&withFollowee, nil).Times(1)
		i.UserRW.EXPECT().Save(gomock.Any()).Return(nil).Times(1)

		retProfile, err := i.GetUCHandler().ProfileUpdateFollow(withFollowee.Name, testData.User("jane").Name, false)
		assert.NoError(t, err)
		assert.Equal(t, testData.User("rick"), *retProfile)
	})
}
