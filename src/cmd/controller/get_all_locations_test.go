package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/cmd/controller"
	"github.com/sanrinconr/storj-images/src/cmd/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	defaultID  = "testID"
	defaultURL = "testURL"
)

func TestEndpoint_Success(t *testing.T) {
	g, err := controller.NewGetAllLocations(defaultGetterMock())
	assert.Nil(t, err)
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	want := minify(t, goldenFile(t, "get_all_locations.golden"))
	err = g(ctx)
	assert.Nil(t, err)

	got := res.Body.Bytes()

	assert.Equal(t, want, got)
}

func TestEndpoint_Fail(t *testing.T) {
	wantErr := errors.New("broken getter")
	var mock mocks.GetterMock = func(ctx context.Context) ([]domain.Location, error) {
		return nil, wantErr
	}
	g, err := controller.NewGetAllLocations(mock)
	assert.Nil(t, err)
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	want := controller.NewError(http.StatusInternalServerError, wantErr)

	got := g(ctx)

	assert.Equal(t, want, got)
}

func TestNewGetAllLocations_fail(t *testing.T) {
	want := errors.New("missing getter in controller get all images")

	_, got := controller.NewGetAllLocations(nil)

	assert.Equal(t, want, got)
}

func defaultGetterMock() mocks.GetterMock {
	return func(ctx context.Context) ([]domain.Location, error) {
		return []domain.Location{
			{
				ID:  defaultID,
				URL: defaultURL,
			},
		}, nil
	}
}

func goldenFile(t *testing.T, fileName string) []byte {
	const location = "testdata/%s"

	f, err := os.Open(fmt.Sprintf(location, fileName))
	assert.Nil(t, err)
	body, err := io.ReadAll(f)
	assert.Nil(t, err)

	return body
}

func minify(t *testing.T, jsonB []byte) []byte {
	buff := &bytes.Buffer{}
	err := json.Compact(buff, jsonB)
	assert.Nil(t, err)

	b, err := io.ReadAll(buff)
	assert.Nil(t, err)

	return b
}
