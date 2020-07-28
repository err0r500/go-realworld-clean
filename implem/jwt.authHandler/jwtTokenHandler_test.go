package jwt_test

import (
	"context"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/stretchr/testify/assert"
)

func TestUserToken_happyCase(t *testing.T) {
	testUserName := "userName"

	tH := jwt.New("theJWTsalt")
	ctx := context.Background()
	token, ok := tH.GenUserToken(ctx, testUserName)
	assert.True(t, ok)

	userName, ok := tH.GetUserName(ctx, token)
	assert.True(t, ok)
	assert.Equal(t, testUserName, userName)
}

func TestUserToken_GenToken_fails(t *testing.T) {
	tH := jwt.New("theJWTsalt")
	token, err := tH.GenUserToken(context.Background(), "")
	assert.False(t, err)
	assert.Equal(t, "", token)
}

func TestUserToken_GetuserName_fails(t *testing.T) {
	tH := jwt.New("theJWTsalt")
	ctx := context.Background()

	token, err := tH.GenUserToken(ctx, "userName")
	assert.True(t, err)

	t.Run("otherSalt", func(t *testing.T) {
		tH2 := jwt.New("otherSalt")
		userName, err := tH2.GetUserName(ctx, token)
		assert.False(t, err)
		assert.Equal(t, "", userName)
	})
}
