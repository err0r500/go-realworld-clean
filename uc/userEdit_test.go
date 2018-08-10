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

func TestEditUser_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedUser := testData.User("rick")
	expectedUser.Email = testData.User("jane").Email
	rick := testData.User("rick")
	jane := testData.User("jane")

	urw := func(i *mock.MockedInteractor) {
		i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
	}

	i := mock.NewMockedInteractor(mockCtrl)
	i.UserRW.EXPECT().GetByName(rick.Password).Return(&rick, nil).AnyTimes()
	urw(&i)

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

func TestInteractor_UserEdit_fails(t *testing.T) {
	rick := testData.User("rick")

	invalidCalls := map[string]mock.TestFunc{
		"error return on uRW.GetByName": func(i *mock.MockedInteractor) {
			i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, errors.New(""))
		},
		"nil, nil return on uRW.GetByName": func(i *mock.MockedInteractor) {
			i.UserRW.EXPECT().GetByName(gomock.Any()).Return(nil, nil)
		},
		"uRW.GetByName returns wrong name": func(i *mock.MockedInteractor) {
			i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&domain.User{Name: "hi there"}, nil)
		},
		"user not validated": func(i *mock.MockedInteractor) {
			i.UserValidator.EXPECT().CheckUser(gomock.Any()).Return(errors.New(""))
		},
		"failed to save the user": func(i *mock.MockedInteractor) {
			i.UserRW.EXPECT().Save(gomock.Any()).Return(errors.New(""))
		},
		"failed to gen token": func(i *mock.MockedInteractor) {
			i.AuthHandler.EXPECT().GenUserToken(gomock.Any()).Return("", errors.New("")).AnyTimes()
		},
	}

	// same as the happy case but with any parameter and called any number of times (including 0)
	validCalls := func(i *mock.MockedInteractor) {
		i.UserRW.EXPECT().GetByName(gomock.Any()).Return(&rick, nil).AnyTimes()
		i.UserValidator.EXPECT().CheckUser(gomock.Any()).Return(nil).AnyTimes()
		i.UserRW.EXPECT().Save(gomock.Any()).Return(nil).AnyTimes()
		i.AuthHandler.EXPECT().GenUserToken(gomock.Any()).Return("token", nil).AnyTimes()
	}

	{ // just to be sure, validCalls don't send any error
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		i := mock.NewMockedInteractor(mockCtrl)
		validCalls(&i)
		_, _, err := i.GetUCHandler().UserEdit("rick", nil)
		assert.NoError(t, err)
	}

	for testName, invalidCall := range invalidCalls {
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			i := mock.NewMockedInteractor(mockCtrl)
			invalidCall(&i) // put the tested call first (important)
			validCalls(&i)  // then fill the gaps with valid calls

			_, _, err := i.GetUCHandler().UserEdit("rick", nil)
			assert.Error(t, err)
		})
	}
}
