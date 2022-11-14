package upload_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/upload"
	"github.com/stretchr/testify/assert"
)

type inserterMock func(context.Context, domain.Image, string) error

func (i inserterMock) Insert(ctx context.Context, img domain.Image, ext string) error {
	if i == nil {
		return errors.New("default mock error")
	}

	return i(ctx, img, ext)
}

func defaultInserterMock() inserterMock {
	return func(ctx context.Context, i domain.Image, s string) error {
		return nil
	}
}

func TestUpload_Success(t *testing.T) {
	a, err := upload.NewAddImage(
		defaultInserterMock(),
		[]string{".jpg"},
		defaultTestTimer(),
	)
	assert.Nil(t, err)
	img := domain.Image{
		Raw:  goldenFile(t, "tree.jpg"),
		Name: "tree.jpg",
	}

	got := a.Upload(context.Background(), img)

	assert.Nil(t, got)
}

func TestUpload_Fail(t *testing.T) {
	t.Run("invalid_format", func(t *testing.T) {
		const allowedExtension = ".raw"
		a, err := upload.NewAddImage(
			defaultInserterMock(),
			[]string{allowedExtension},
			defaultTestTimer(),
		)
		assert.Nil(t, err)
		img := domain.Image{
			Raw:  goldenFile(t, "tree.jpg"),
			Name: "tree.jpg",
		}
		want := domain.InvalidFormatImageError(errors.New("invalid image format"))

		got := a.Upload(context.Background(), img)

		assert.Equal(t, want, got)
	})

	t.Run("broken_inserter", func(t *testing.T) {
		want := errors.New("broken inserter")
		var inserter inserterMock = func(ctx context.Context, i domain.Image, s string) error {
			return want
		}
		a, err := upload.NewAddImage(
			inserter,
			[]string{".jpg"},
			defaultTestTimer(),
		)
		assert.Nil(t, err)
		img := domain.Image{
			Raw:  goldenFile(t, "tree.jpg"),
			Name: "tree.jpg",
		}

		got := a.Upload(context.Background(), img)

		assert.Equal(t, want, got)
	})
}

func TestNew_Fail(t *testing.T) {
	t.Run("inserter", func(t *testing.T) {
		want := errors.New("missing inserter in add image func")

		_, got := upload.NewAddImage(nil, []string{".jpg"}, defaultTestTimer())

		assert.Equal(t, want, got)
	})

	t.Run("formats", func(t *testing.T) {
		want := errors.New("missing formats in add image func")

		_, got := upload.NewAddImage(defaultInserterMock(), []string{}, defaultTestTimer())

		assert.Equal(t, want, got)
	})
}
