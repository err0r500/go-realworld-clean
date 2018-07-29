package domain_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	user := testData.User("rick")
	jane := testData.User("jane")
	emptyString := ""
	newBio := "new biography"

	t.Run("nil pointers should leave struct intact", func(t *testing.T) {
		domain.UpdateUser(&user,
			domain.SetUserEmail(nil),
			domain.SetUserName(nil),
			domain.SetUserBio(nil),
			domain.SetUserImageLink(nil),
			domain.SetUserPassword(nil),
		)

		assert.Equal(t, testData.User("rick"), user)
	})

	t.Run("update all structs", func(t *testing.T) {
		domain.UpdateUser(&user,
			domain.SetUserEmail(&jane.Email),
			domain.SetUserName(&jane.Name),
			domain.SetUserBio(&newBio),
			domain.SetUserImageLink(testData.User("jane").ImageLink),
			domain.SetUserPassword(&jane.Password),
		)

		assert.Equal(t, jane.Email, user.Email)
		assert.Equal(t, jane.Name, user.Name)
		assert.Equal(t, newBio, *user.Bio)
		assert.Equal(t, jane.ImageLink, user.ImageLink)
		assert.Equal(t, jane.Password, user.Password)
	})

	t.Run("empty bio string deletes bio property", func(t *testing.T) {
		domain.UpdateUser(&user,
			domain.SetUserBio(&emptyString),       // should set to nil
			domain.SetUserImageLink(&emptyString), // should set to nil
		)
		assert.Nil(t, user.Bio)
		assert.Nil(t, user.ImageLink)
	})

}

func TestFollows(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		userFollows := []string{"a", "b", "c", "d", "f", "e"}
		user := domain.User{
			FollowIDs: userFollows,
		}

		assert.True(t, user.Follows("e"))
		assert.False(t, user.Follows("e "))
		assert.False(t, user.Follows(""))
	})
	t.Run("standard", func(t *testing.T) {
		user := domain.User{
			FollowIDs: nil,
		}
		assert.False(t, user.Follows("a"))
	})
}

func TestUser_UpdateFollowees(t *testing.T) {

	t.Run("add", func(t *testing.T) {
		user := domain.User{
			FollowIDs: []string{"a", "b", "c"},
		}
		user.UpdateFollowees("d", true)
		user.UpdateFollowees("e", true)
		user.UpdateFollowees("f", true)
		assert.True(t, user.Follows("d"))
		assert.True(t, user.Follows("e"))
		assert.True(t, user.Follows("f"))
	})

	t.Run("remove", func(t *testing.T) {
		user := domain.User{
			FollowIDs: []string{"a", "b", "c", "d", "e", "f"},
		}
		user.UpdateFollowees("d", false)
		user.UpdateFollowees("e", false)
		user.UpdateFollowees("f", false)

		assert.False(t, user.Follows("d"))
		assert.False(t, user.Follows("e"))
		assert.False(t, user.Follows("f"))
	})

	t.Run("remove not present", func(t *testing.T) {
		user := domain.User{
			FollowIDs: []string{"a", "b", "c"},
		}
		user.UpdateFollowees("d", false)
		assert.False(t, user.Follows("d"))
	})

	t.Run("remove while empty", func(t *testing.T) {
		user := domain.User{
			FollowIDs: nil,
		}
		user.UpdateFollowees("d", false)
		assert.False(t, user.Follows("d"))
	})

}
