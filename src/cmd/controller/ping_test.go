package controller_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/controller"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctl := controller.NewPing()

	got := ctl(c)

	assert.Nil(t, got)
}
