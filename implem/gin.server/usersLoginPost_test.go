package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gopkg.in/h2non/baloo.v3"
)

var userLoginPostPath = "/api/users/login"

var expectedRespSchema = `{
	"type": "object",
  "additionalProperties": false,
  "properties": {
		"user": {
			"type": "object",
			"additionalProperties": false,
			"properties": {
				"email": {
					"type": "string",
					"format": "email"
				},
				"token": {
					"type": "string",
 					"minLength": 1
				},
				"username": {
					"type": "string",
 					"minLength": 1
				},
				"bio": {
					"anyOf": [
						{ "type": "string" },
						{ "type": "null" }
					]
				},
				"image": {
					"anyOf": [
						{ "type": "string" },
						{ "type": "null" }
					]
				}
			}
		}
	}
}`

func TestUserLoginPost_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jane := testData.User("jane")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		UserLogin(jane.Email, jane.Password).
		Return(&jane, "authToken", nil).
		Times(1)

	jwtHandler := jwt.NewTokenHandler("mySalt")
	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Post(userLoginPostPath).
		BodyString(`
		{
  			"user": {
				"email": "` + testData.User("jane").Email + `",
    			"password": "` + testData.User("jane").Password + `"
  			}
		}`).
		Expect(t).
		Status(http.StatusOK).
		JSONSchema(expectedRespSchema).
		Done()
}
