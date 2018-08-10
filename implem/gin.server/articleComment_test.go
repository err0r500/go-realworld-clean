package server_test

import (
	"net/http/httptest"
	"testing"

	"strconv"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var articleCommentPath = "/api/articles/" + testData.Article("jane").Slug + "/comments"

func TestArticleCommentGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedComment := testData.Article("").Comments
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		CommentsGet(
			testData.Article("jane").Slug,
		).
		Return(expectedComment, nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	t.Run("happyCase", func(t *testing.T) {
		baloo.New(ts.URL).
			Get(articleCommentPath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			Expect(t).
			Status(200).
			JSONSchema(testData.CommentsMultipleResponse).
			Done()
	})
}

func TestArticleCommentPost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedComment := testData.Article("").Comments[0]
	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		CommentsPost(
			testData.User("jane").Name,
			testData.Article("jane").Slug,
			"a comment",
		).
		Return(&expectedComment, nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	validReq := `{
  "comment": {
    "body": "a comment"
  }
}`

	t.Run("happyCase", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(articleCommentPath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			BodyString(validReq).
			Expect(t).
			Status(201).
			JSONSchema(testData.CommentsSimgleResponse).
			Done()
	})

	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(articleCommentPath).
			BodyString(validReq).
			Expect(t).
			Status(401).
			Done()

	})

	t.Run("no body", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(articleCommentPath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			Expect(t).
			Status(400).
			Done()
	})
}

func TestArticleCommentDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		CommentsDelete(
			testData.User("jane").Name,
			testData.Article("jane").Slug,
			testData.Article("jane").Comments[0].ID,
		).
		Return(nil).
		Times(1)

	jwtHandler := jwt.New("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	t.Run("happyCase", func(t *testing.T) {
		baloo.New(ts.URL).
			Delete(articleCommentPath+"/"+strconv.Itoa(testData.Article("jane").Comments[0].ID)).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			Expect(t).
			Status(200).
			Done()
	})

	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Delete(articleCommentPath + "/" + strconv.Itoa(testData.Article("jane").Comments[0].ID)).
			Expect(t).
			Status(401).
			Done()
	})
}
