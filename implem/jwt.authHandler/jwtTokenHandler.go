package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/err0r500/go-realworld-clean/uc"
)

var tokenTimeToLive = time.Hour * 2

// tokenHandler handles JWT related request, implementing uc.AuthHandler interface
type tokenHandler struct {
	salt []byte
}

func New(salt string) uc.AuthHandler {
	return tokenHandler{
		salt: []byte(salt),
	}
}

//GenToken (uc.Admin) : returns a signed token for an admin
func (tH tokenHandler) GenUserToken(userName string) (string, error) {
	if userName == "" {
		return "", errors.New("can't generate token for empty user")
	}

	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, newUserClaims(userName, tokenTimeToLive)).
		SignedString(tH.salt)
}

func (tH tokenHandler) GetUserName(inToken string) (userName string, err error) {
	token, err := jwt.ParseWithClaims(
		inToken,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tH.salt), nil
		},
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims.Name, nil
	}

	return "", errors.New("problem with jwt")
}

type userClaims struct {
	Name string
	jwt.StandardClaims
}

// newUserClaims : constructor of userClaims
func newUserClaims(name string, ttl time.Duration) *userClaims {
	return &userClaims{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Issuer:    "real-worl-demo-backend",
		},
	}
}
