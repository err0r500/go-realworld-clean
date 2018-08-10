package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gopkg.in/h2non/baloo.v3"
)

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
