package domain_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	user := testData.User("rick")

	domain.UpdateUser(&user,
		domain.SetUserEmail(nil),
		domain.SetUserName(nil),
		domain.SetUserBio(nil),
		domain.SetUserImageLink(nil),
		domain.SetUserPassword(nil),
	)

	assert.Equal(t, testData.User("rick"), user)

	emptyString := ""
	jane := testData.User("jane")
	domain.UpdateUser(&user,
		domain.SetUserEmail(&jane.Email),
		domain.SetUserName(&jane.Name),
		domain.SetUserBio(&emptyString), // should set to nil
		domain.SetUserImageLink(testData.User("jane").ImageLink),
		domain.SetUserPassword(&jane.Password),
	)

	assert.Equal(t, testData.User("jane").Email, user.Email)
	assert.Equal(t, testData.User("jane").Name, user.Name)
	assert.Nil(t, user.Bio)
	assert.Equal(t, testData.User("jane").ImageLink, user.ImageLink)
	assert.Equal(t, testData.User("jane").Password, user.Password)
}

func TestFollows(t *testing.T) {
	userFollows := []string{"a", "b", "c", "d", "f", "e"}
	user := domain.User{
		FollowIDs: userFollows,
	}

	assert.True(t, user.Follows("e"))
	assert.False(t, user.Follows("e "))
	assert.False(t, user.Follows(""))
}
