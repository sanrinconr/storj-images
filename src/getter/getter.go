// Package getter obtain saved images from the cloud.
package getter

import (
	"context"
	"fmt"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/getter/internal"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type (
	metadata interface {
		//nolint:godox // to think later.
		// TODO: this abstraction make sense?, need a bson, a mongo package.
		GetAll(ctx context.Context, query, projections bson.M) (*mongodriver.Cursor, error)
	}

	object interface {
		GetShareableLink(context.Context, string) (string, error)
	}
)

// Getter has a metada to obtain location and object is the object storage cloud.
type Getter struct {
	metadata metadata
	object   object
}

// New create a getter object.
func New(m metadata, o object) (Getter, error) {
	g := Getter{
		metadata: m,
		object:   o,
	}

	return g, g.validate()
}

// All obtain all images saved in metadata, with this make a query into object storage.
func (g Getter) All(ctx context.Context) ([]domain.Location, error) {
	c, err := g.metadata.GetAll(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	res := make([]domain.Location, c.RemainingBatchLength())

	for i := 0; c.Next(ctx); i++ {
		var img internal.Image
		if err := c.Decode(&img); err != nil {
			return nil, err
		}

		loc := domain.Location{}
		loc.ID = img.ObjectStorageKey
		loc.Name = img.Name
		loc.CreatedAt = img.CreatedAt

		loc.URL, err = g.object.GetShareableLink(ctx, img.ObjectStorageKey)
		if err != nil {
			return nil, err
		}

		res[i] = loc
	}

	return res, nil
}

func (g Getter) validate() error {
	const missingDependency = "missing %s dependency in getter"

	if g.metadata == nil {
		return fmt.Errorf(missingDependency, "metadata storage")
	}

	if g.object == nil {
		return fmt.Errorf(missingDependency, "object storage")
	}

	return nil
}
