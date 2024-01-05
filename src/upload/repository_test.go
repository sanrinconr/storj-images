package upload_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/upload"
	"github.com/stretchr/testify/assert"
)

type objectMock func(ctx context.Context, key string, val []byte) error

func (m objectMock) Insert(ctx context.Context, key string, val []byte) error {
	if m == nil {
		return errors.New("default mock error")
	}

	return m(ctx, key, val)
}

func defaultObjectMock() objectMock {
	return func(ctx context.Context, key string, val []byte) error {
		return nil
	}
}

func TestInsert_Success(t *testing.T) {
	var gotCallsObject int
	var objectMock objectMock = func(ctx context.Context, key string, val []byte) error {
		gotCallsObject++

		return nil
	}
	r, err := upload.NewRepository(objectMock, defaultTestTimer())
	assert.Nil(t, err)
	wantCallsObject := 1

	got := r.Insert(context.Background(), domain.Image{}, ".test")

	assert.Nil(t, got)
	assert.Equal(t, wantCallsObject, gotCallsObject)
}

func TestInsert_Fails(t *testing.T) {
	t.Run("broken_object_storage", func(t *testing.T) {
		want := errors.New("broken storage")

		var objectMock objectMock = func(ctx context.Context, key string, val []byte) error {
			return want
		}
		r, err := upload.NewRepository(objectMock, defaultTestTimer())
		assert.Nil(t, err)

		got := r.Insert(context.Background(), domain.Image{}, ".test")

		assert.Equal(t, want, got)
	})
}

func TestNew_Fails(t *testing.T) {
	const template = "missing dependency %s in image repository"
	t.Run("object", func(t *testing.T) {
		want := fmt.Errorf(template, "object infra")

		_, got := upload.NewRepository(nil, defaultTestTimer())

		assert.Equal(t, want, got)
	})
	t.Run("timer", func(t *testing.T) {
		want := fmt.Errorf(template, "timer")

		_, got := upload.NewRepository(defaultObjectMock(), nil)

		assert.Equal(t, want, got)
	})
}
