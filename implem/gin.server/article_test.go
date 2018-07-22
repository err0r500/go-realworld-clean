package server_test

import (
	"testing"

	"net/http/httptest"

	"net/http"

	server "github.com/err0r500/go-realworld-clean/implem/gin.server"
	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	"github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

const artPath = "/api/articles"

func TestRouterHandler_articlePost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedArticle := testData.Article("jane")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ArticlePost(testData.User("jane").Name, gomock.Any()).
		Return(&expectedArticle, nil).
		Times(1)

	jwtHandler := jwt.NewTokenHandler("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	reqBody := `{
  "article": {
    "title": "` + expectedArticle.Title + `",
    "description": "` + expectedArticle.Description + `",
    "body": "` + expectedArticle.Body + `",
    "tagList": [
      "` + expectedArticle.TagList[0] + `"
    ]
  }
}`
	t.Run("happycase", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(artPath).
			AddHeader("Authorization", authToken).
			BodyString(reqBody).Expect(t).
			Status(http.StatusCreated).
			Done()
	})
	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(artPath).
			BodyString(reqBody).Expect(t).
			Status(http.StatusUnauthorized).
			Done()
	})
	t.Run("no body", func(t *testing.T) {
		baloo.New(ts.URL).
			Post(artPath).
			AddHeader("Authorization", authToken).
			Expect(t).
			Status(http.StatusBadRequest).
			Done()
	})
}

func TestRouterHandler_articlePut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedArticle := testData.Article("jane")
	jwtHandler := jwt.NewTokenHandler("mySalt")

	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ArticlePut(testData.User("jane").Name, expectedArticle.Slug, gomock.Any()).
		Return(&expectedArticle, nil).
		Times(1)

	gE := gin.Default()
	server.NewRouterWithLogger(ucHandler, jwtHandler, uc.SimpleLogger{}).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	reqBody := `{
  "article": {
    "title": "` + expectedArticle.Title + `",
    "description": "` + expectedArticle.Description + `",
    "body": "` + expectedArticle.Body + `",
    "tagList": [
      "` + expectedArticle.TagList[0] + `"
    ]
  }
}`

	t.Run("happycase", func(t *testing.T) {
		baloo.New(ts.URL).
			Put(artPath+"/"+expectedArticle.Slug).
			AddHeader("Authorization", authToken).
			BodyString(reqBody).
			Expect(t).
			Status(http.StatusOK).
			Done()
	})

	t.Run("no auth", func(t *testing.T) {
		baloo.New(ts.URL).
			Put(artPath + "/" + expectedArticle.Slug).
			BodyString(reqBody).
			Expect(t).
			Status(http.StatusUnauthorized).
			Done()
	})

	t.Run("no body", func(t *testing.T) {
		baloo.New(ts.URL).
			Put(artPath+"/"+expectedArticle.Slug).
			AddHeader("Authorization", authToken).
			Expect(t).
			Status(http.StatusBadRequest).
			Done()
	})
}

func TestRouterHandler_articleGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedArticle := testData.Article("jane")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ArticleGet(expectedArticle.Slug).
		Return(&expectedArticle, nil).
		Times(1)

	gE := gin.Default()
	server.NewRouterWithLogger(ucHandler, jwt.NewTokenHandler("mySalt"), uc.SimpleLogger{}).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Get(artPath + "/" + expectedArticle.Slug).
		Expect(t).
		Status(http.StatusOK).
		Done()
}

func TestRouterHandler_articleDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedArticle := testData.Article("jane")
	ucHandler := uc.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		ArticleDelete(testData.User("jane").Name, expectedArticle.Slug).
		Return(nil).
		Times(1)

	jwtHandler := jwt.NewTokenHandler("mySalt")

	gE := gin.Default()
	server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	authToken, err := jwtHandler.GenUserToken(testData.User("jane").Name)
	assert.NoError(t, err)

	baloo.New(ts.URL).
		Delete(artPath+"/"+expectedArticle.Slug).
		AddHeader("Authorization", authToken).
		Expect(t).
		Status(http.StatusOK).
		Done()

	baloo.New(ts.URL).
		Delete(artPath + "/" + expectedArticle.Slug).
		Expect(t).
		Status(http.StatusUnauthorized).
		Done()
}
