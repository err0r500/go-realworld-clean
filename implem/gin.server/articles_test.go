package server_test

import (
	"net/http/httptest"
	"testing"

	"strconv"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	jwt "github.com/err0r500/go-realworld-clean/implem/jwt.authHandler"
	mock "github.com/err0r500/go-realworld-clean/implem/mock.uc"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/baloo.v3"
)

var articlesFilteredPath = "/api/articles"
var articlesFeedPath = "/api/articles/feed"

func TestArticlesFiltered(t *testing.T) {
	// todo : add test with auth
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	limit := 10
	offset := 2
	author := "jane"
	tag := "tag1"
	fav := "false"

	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		GetArticles("", limit, offset, gomock.Any()).
		Return(nil, domain.ArticleCollection{testData.Article("jane")}, 10, nil).
		Times(1)

	gE := gin.Default()
	server.NewRouter(ucHandler, nil).SetRoutes(gE)
	ts := httptest.NewServer(gE)
	defer ts.Close()

	t.Run("happyCase", func(t *testing.T) {
		baloo.New(ts.URL).
			Get(articlesFilteredPath).
			AddQuery("limit", strconv.Itoa(limit)).
			AddQuery("offset", strconv.Itoa(offset)).
			AddQuery("author", author).
			AddQuery("tag", tag).
			AddQuery("favorited", fav).
			Expect(t).
			Status(200).
			JSONSchema(testData.ArticleMultipleRespDefinition).
			Done()
	})
}

func TestArticlesFeed(t *testing.T) {
	t.Run("happyCase", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		limit := 10
		offset := 2

		jane := testData.User("jane")

		ucHandler := mock.NewMockHandler(mockCtrl)
		ucHandler.EXPECT().
			ArticlesFeed(jane.Name, limit, offset).
			Return(nil, domain.ArticleCollection{testData.Article("jane")}, 10, nil).
			Times(1)

		jwtHandler := jwt.NewTokenHandler("mySalt")

		gE := gin.Default()
		server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
		ts := httptest.NewServer(gE)
		defer ts.Close()

		authToken, err := jwtHandler.GenUserToken(jane.Name)
		assert.NoError(t, err)

		baloo.New(ts.URL).
			Get(articlesFeedPath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			AddQuery("limit", strconv.Itoa(limit)).
			AddQuery("offset", strconv.Itoa(offset)).
			Expect(t).
			Status(200).
			JSONSchema(testData.ArticleMultipleRespDefinition).
			Done()
	})

	t.Run("empty", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		limit := 10
		offset := 2

		jane := testData.User("jane")

		ucHandler := mock.NewMockHandler(mockCtrl)
		ucHandler.EXPECT().
			ArticlesFeed(jane.Name, limit, offset).
			Return(&jane, nil, 0, nil).
			Times(1)

		jwtHandler := jwt.NewTokenHandler("mySalt")

		gE := gin.Default()
		server.NewRouter(ucHandler, jwtHandler).SetRoutes(gE)
		ts := httptest.NewServer(gE)
		defer ts.Close()

		authToken, err := jwtHandler.GenUserToken(jane.Name)
		assert.NoError(t, err)

		baloo.New(ts.URL).
			Get(articlesFeedPath).
			AddHeader("Authorization", testData.TokenPrefix+authToken).
			AddQuery("limit", strconv.Itoa(limit)).
			AddQuery("offset", strconv.Itoa(offset)).
			Expect(t).
			Status(200).
			BodyEquals(`{"articles":[],"articlesCount":0}`).
			Done()
	})
}
