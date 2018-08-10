package server_test

import (
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var profileFollowPostPath = "/api/profiles/" + testData.User("rick").Name + "/follow"

func TestProfileFollowPost_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	rick := testData.User("rick")
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ProfileUpdateFollow(testData.User("jane").Name, rick.Name, true).
		Return(&rick, nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")
	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Post(profileFollowPostPath).
		AddHeader("Authorization", testData.TokenPrefix+authToken).
		Expect(t).
		Status(200).
		JSONSchema(testData.ProfileRespDefinition).
		Done()
}
