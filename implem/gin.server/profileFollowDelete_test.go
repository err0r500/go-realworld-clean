package server_test

import (
	"net/http/httptest"
	"testing"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var profileFollowDeletePath = "/api/profiles/" + testData.User("rick").Name + "/follow"

func TestProfileFollowDelete_happyCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ProfileUpdateFollow(testData.User("jane").Name, testData.User("rick").Name, false).
		Return(&domain.User{}, nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	baloo.New(ts.URL).
		Delete(profileFollowDeletePath).
		AddHeader("Authorization", testData.TokenPrefix+authToken).
		Expect(t).
		Status(200).
		JSONSchema(testData.ProfileRespDefinition).
		Done()
}
