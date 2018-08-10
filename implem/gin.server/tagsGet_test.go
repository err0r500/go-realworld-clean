package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/err0r500/go-realworld-clean/implem/gin.server"
	"github.com/err0r500/go-realworld-clean/implem/uc.mock"
	"github.com/err0r500/go-realworld-clean/testData"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"gopkg.in/h2non/baloo.v3"
)

var tagsPath = "/api/tags"

func TestTagsGet_happyCase(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		tags := []string{"tag1", "tag2"}
		ucHandler := mock.NewMockHandler(mockCtrl)
		ucHandler.EXPECT().
			Tags().
			Return(tags, nil).
			Times(1)

		gE := gin.Default()
		server.NewRouter(ucHandler, nil).SetRoutes(gE)

		ts := httptest.NewServer(gE)
		defer ts.Close()

		baloo.New(ts.URL).
			Get(tagsPath).
			Expect(t).
			Status(http.StatusOK).
			JSONSchema(testData.TagsResponse).
			Done()
	})
	t.Run("empty", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		ucHandler := mock.NewMockHandler(mockCtrl)
		ucHandler.EXPECT().
			Tags().
			Return(nil, nil).
			Times(1)

		gE := gin.Default()
		server.NewRouter(ucHandler, nil).SetRoutes(gE)

		ts := httptest.NewServer(gE)
		defer ts.Close()

		baloo.New(ts.URL).
			Get(tagsPath).
			Expect(t).
			Status(http.StatusOK).
			BodyEquals(`{"tags":[]}`).
			Done()
	})

}

func TestTagsGet_fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ucHandler := mock.NewMockHandler(mockCtrl)
	ucHandler.EXPECT().
		Tags().
		Return(nil, errors.New("")).
		Times(1)

	gE := gin.Default()
	server.NewRouterWithLogger(ucHandler, nil, mock.SimpleLogger{}).SetRoutes(gE)

	ts := httptest.NewServer(gE)
	defer ts.Close()

	baloo.New(ts.URL).
		Get(tagsPath).
		Expect(t).
		Status(http.StatusUnprocessableEntity).
		Done()
}
