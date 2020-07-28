package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"

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

// GenUserToken (uc.Admin) : returns a signed token for an admin
func (tH tokenHandler) GenUserToken(ctx context.Context, userName string) (string, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "auth:gen_token")
	defer span.Finish()

	if userName == "" {
		return "", false
	}

	token, err := jwt.
		NewWithClaims(jwt.SigningMethodHS256, newUserClaims(userName, tokenTimeToLive)).
		SignedString(tH.salt)
	if err != nil {
		return "", false
	}

	return token, true
}

func (tH tokenHandler) GetUserName(ctx context.Context, inToken string) (string, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "auth:get_username")
	defer span.Finish()

	token, err := jwt.ParseWithClaims(
		inToken,
		&userClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tH.salt), nil
		},
	)
	if err != nil {
		span.LogFields(log.Error(errors.New("parsing failed")))
		return "", false
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		return claims.Name, true
	}

	span.LogFields(log.Error(errors.New("not userClaims")))

	return "", false
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
