package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	server "github.com/err0r500/go-realworld-clean/implem/gin.server"
	logger "github.com/err0r500/go-realworld-clean/implem/logrus.logger"
	mock "github.com/err0r500/go-realworld-clean/implem/uc.mock"

	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var userGetPath = "/api/user"

func TestUserGet_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jane := testData.User("jane")
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		UserGet(jane.Name).
		Return(&jane, "authToken", nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")
	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Get(userGetPath).
		AddHeader("Authorization", testData.TokenPrefix+authToken).
		Expect(t).
		Status(http.StatusOK).
		JSONSchema(testData.UserRespDefinition).
		Done()
}

var userPutPath = "/api/user"

func TestUserPut_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jane := testData.User("jane")
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		UserEdit(testData.User("rick").Name, gomock.Any()).
		Return(&jane, "token", nil).
		Times(1)

	gE := gin.Default()
	jwtHandler := jwt.New("mySalt")
	router := server.NewRouter(ucHandler, jwtHandler)
	router.Logger = logger.SimpleLogger{}
	router.SetRoutes(gE)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	token, err := jwtHandler.GenUserToken(testData.User("rick").Name)
	assert.NoError(t, err)
	baloo.New(ts.URL).
		Put(userPutPath).
		AddHeader("Authorization", testData.TokenPrefix+token).
		BodyString(`{
  			"user": {
				"bio": "` + testData.User("rick").Email + `",
				"name": "` + testData.User("rick").Name + `"
  			}
		}`).
		Expect(t).
		Status(http.StatusOK).
		JSONSchema(testData.UserRespDefinition).
		Done()
}

var userPostPath = "/api/users"

func TestUserPost_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	jane := testData.User("jane")
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		UserCreate(jane.Name, jane.Email, jane.Password).
		Return(&jane, "authToken", nil).
		Times(1)

	gE := gin.Default()
	server.NewRouter(ucHandler, jwt.New("mySalt")).SetRoutes(gE)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Post(userPostPath).
		BodyString(`
		{
  			"user": {
    			"username": "` + testData.User("jane").Name + `",
				"email": "` + testData.User("jane").Email + `",
    			"password": "` + testData.User("jane").Password + `"
  			}
		}`).
		Expect(t).
		JSONSchema(testData.UserRespDefinition).
		Status(http.StatusCreated).
		Done()
}
