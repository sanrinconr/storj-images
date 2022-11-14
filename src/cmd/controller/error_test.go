package controller_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sanrinconr/storj-images/src/cmd/controller"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {

	t.Run("fe", func(t *testing.T) {
		err := controller.NewError(http.StatusOK, errors.New("pruebas"))
		want := "pruebas"

		got := err.Error()

		assert.Equal(t, want, got)
	})
}
