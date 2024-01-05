package getter_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/getter"
	"github.com/stretchr/testify/assert"
)

var (
	defaultID           = "123"
	defaultName         = "tree.png"
	defaultStorageKey   = "123.png"
	defaultShareableURL = "test.com/123.jpg"
)

var errDefaultMock = errors.New("default mock error")

type objectMock struct {
	getShareableLinkMock func(context.Context, string) (string, error)
	obtainAllKeys        func(ctx context.Context) ([]string, error)
}

func (o objectMock) ObtainAllKeys(ctx context.Context) ([]string, error) {
	if o.obtainAllKeys == nil {
		return nil, errDefaultMock
	}

	return o.obtainAllKeys(ctx)
}

func (o objectMock) GetShareableLink(ctx context.Context, key string) (string, error) {
	if o.getShareableLinkMock == nil {
		return "", errDefaultMock
	}

	return o.getShareableLinkMock(ctx, key)
}

func TestAll_Success(t *testing.T) {
	g, err := getter.New(defaultObjectMock())
	assert.Nil(t, err)
	want := []domain.Location{
		{
			ID:  defaultID,
			URL: defaultShareableURL,
		},
	}

	got, err := g.All(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestAll_Fails(t *testing.T) {
	t.Run("object_storage_shareable_link_error", func(t *testing.T) {
		var m objectMock
		g, err := getter.New(m)
		assert.Nil(t, err)
		want := errDefaultMock

		_, got := g.All(context.Background())

		assert.Equal(t, got, want)
	})
}

func TestNew_fail(t *testing.T) {
	t.Run("nil_object", func(t *testing.T) {
		want := errors.New("missing object storage dependency in getter")

		_, got := getter.New(nil)

		assert.Equal(t, want, got)
	})
}

func defaultObjectMock() objectMock {
	return objectMock{
		getShareableLinkMock: func(ctx context.Context, s string) (string, error) {
			return defaultShareableURL, nil
		},
		obtainAllKeys: func(ctx context.Context) ([]string, error) {
			return []string{defaultID}, nil
		},
	}
}
