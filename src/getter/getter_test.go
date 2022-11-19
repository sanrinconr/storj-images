package getter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sanrinconr/storj-images/src/getter"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type metadataMock func(context.Context, bson.M, bson.M) (*mongo.Cursor, error)

func (m metadataMock) GetAll(ctx context.Context, query bson.M, projection bson.M) (*mongo.Cursor, error) {
	if m == nil {
		return nil, errors.New("default mock error")
	}

	return m(ctx, query, projection)
}

type objectMock func(context.Context, string) (string, error)

func (o objectMock) GetShareableLink(ctx context.Context, key string) (string, error) {
	if o == nil {
		return "", errors.New("default error mock")
	}

	return o(ctx, key)
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
