package getter_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/getter"
	"github.com/sanrinconr/storj-images/src/mongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	defaultID           = "123"
	defaultName         = "tree.png"
	defaultStorageKey   = "123.png"
	defaultShareableURL = "test.com/123.jpg"
	defaultTime         = time.Date(2019, time.August, 12, 0, 0, 0, 0, time.UTC)
)

var errDefaultMock = errors.New("default mock error")

type metadataMock func(context.Context, bson.M, bson.M) ([]mongo.Document, error)

func (m metadataMock) GetAll(ctx context.Context, query bson.M, projection bson.M) ([]mongo.Document, error) {
	if m == nil {
		return nil, errDefaultMock
	}

	return m(ctx, query, projection)
}

type objectMock func(context.Context, string) (string, error)

func (o objectMock) GetShareableLink(ctx context.Context, key string) (string, error) {
	if o == nil {
		return "", errDefaultMock
	}

	return o(ctx, key)
}

func TestAll_Success(t *testing.T) {
	g, err := getter.New(defaultMetadataMock(), defaultObjectMock())
	assert.Nil(t, err)
	want := []domain.Location{
		{
			ID:        defaultID,
			Name:      defaultName,
			URL:       defaultShareableURL,
			CreatedAt: defaultTime,
		},
	}

	got, err := g.All(context.Background())

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestAll_Fails(t *testing.T) {
	t.Run("metadata_error", func(t *testing.T) {
		var m metadataMock
		g, err := getter.New(m, defaultObjectMock())
		assert.Nil(t, err)
		want := errDefaultMock

		_, got := g.All(context.Background())

		assert.Equal(t, got, want)
	})

	t.Run("object_storage_shareable_link_error", func(t *testing.T) {
		var m objectMock
		g, err := getter.New(defaultMetadataMock(), m)
		assert.Nil(t, err)
		want := errDefaultMock

		_, got := g.All(context.Background())

		assert.Equal(t, got, want)
	})
}

func TestNew_fail(t *testing.T) {
	t.Run("nil_metadata", func(t *testing.T) {
		want := errors.New("missing metadata storage dependency in getter")
		var mock objectMock

		_, got := getter.New(nil, mock)

		assert.Equal(t, want, got)
	})

	t.Run("nil_object", func(t *testing.T) {
		want := errors.New("missing object storage dependency in getter")
		var mock metadataMock

		_, got := getter.New(mock, nil)

		assert.Equal(t, want, got)
	})
}

func defaultMetadataMock() metadataMock {
	return func(ctx context.Context, m1, m2 bson.M) ([]mongo.Document, error) {
		return []mongo.Document{
			{
				ID:               defaultID,
				Name:             defaultName,
				ObjectStorageKey: defaultStorageKey,
				CreatedAt:        time.Date(2019, time.August, 12, 0, 0, 0, 0, time.UTC),
			},
		}, nil
	}
}

func defaultObjectMock() objectMock {
	return func(ctx context.Context, s string) (string, error) {
		return defaultShareableURL, nil
	}
}
