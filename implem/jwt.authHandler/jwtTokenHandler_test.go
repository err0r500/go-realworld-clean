package jwt_test

import (
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/stretchr/testify/assert"
)

func TestUserToken_happyCase(t *testing.T) {
	testUserName := "userName"

	tH := jwt.New("theJWTsalt")
	token, err := tH.GenUserToken(testUserName)
	assert.NoError(t, err)

	userName, err := tH.GetUserName(token)
	assert.NoError(t, err)
	assert.Equal(t, testUserName, userName)
}

func TestUserToken_GenToken_fails(t *testing.T) {
	tH := jwt.New("theJWTsalt")
	token, err := tH.GenUserToken("")
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestUserToken_GetuserName_fails(t *testing.T) {
	tH := jwt.New("theJWTsalt")
	token, err := tH.GenUserToken("userName")
	assert.NoError(t, err)

	t.Run("otherSalt", func(t *testing.T) {
		tH2 := jwt.New("otherSalt")
		userName, err := tH2.GetUserName(token)
		assert.Error(t, err)
		assert.Equal(t, "", userName)
	})
}
