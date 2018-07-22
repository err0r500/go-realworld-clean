package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var profileGetPath = "/api/profiles/" + testData.User("rick").Name

func TestProfileGet_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	janeFollowingRick := testData.Profile("janeFollowingRick")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ProfileGet(testData.User("rick").Name).
		Return(&janeFollowingRick, nil).
		Times(1)

	gE := gin.Default()
	server.NewRouter(ucHandler, jwt.NewTokenHandler("mySalt")).SetRoutes(gE)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	rick := testData.User("rick")
	baloo.New(ts.URL).
		Get(profileGetPath).
		Expect(t).
		Status(200).
		AssertFunc(func(response *http.Response, request *http.Request) error {
			body, err := ioutil.ReadAll(response.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"profile":{"username":"`+rick.Name+`","bio":"`+*rick.Bio+`", "following": true}}`, string(body))
			return nil
		}).
		Done()
}
