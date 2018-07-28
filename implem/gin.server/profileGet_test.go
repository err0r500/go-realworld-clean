package server_test

import (
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

	baloo.New(ts.URL).
		Get(profileGetPath).
		Expect(t).
		Status(200).
		JSONSchema(testData.ProfileRespDefinition).
		Done()
}
