package controller_test

import (
	"errors"
	"testing"

	"github.com/sanrinconr/storj-images/src/cmd/controller"
	"github.com/stretchr/testify/assert"
)

func TestNewAddImage_fails(t *testing.T) {
	want := errors.New("missing uploader in add image controller")

	_, got := controller.NewAddImage(nil)

	assert.Equal(t, want, got)
}

//nolint:godox // to keep in mind latter.
// TODO: make test of multipart, i try, but is hard, need more time to think.
