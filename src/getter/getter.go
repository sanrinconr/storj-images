// Package getter obtain saved images from the cloud.
package getter

import (
	"context"
	"fmt"

	domain "github.com/sanrinconr/storj-images/src"
)

type object interface {
	GetShareableLink(context.Context, string) (string, error)
	ObtainAllKeys(ctx context.Context) ([]string, error)
}

// Getter has a metada to obtain location and object is the object storage cloud.
type Getter struct {
	object object
}

// New create a getter object.
func New(o object) (Getter, error) {
	g := Getter{
		object: o,
	}

	return g, g.validate()
}

// All obtain all images saved in metadata, with this make a query into object storage.
func (g Getter) All(ctx context.Context) ([]domain.Location, error) {
	docs, err := g.object.ObtainAllKeys(ctx)
	if err != nil {
		return nil, err
	}

	locs := make([]domain.Location, 0, len(docs))

	for _, doc := range docs {
		loc := domain.Location{
			ID: doc,
		}

		loc.URL, err = g.object.GetShareableLink(ctx, doc)
		if err != nil {
			return nil, err
		}

		locs = append(locs, loc)
	}

	return locs, nil
}

func (g Getter) validate() error {
	const missingDependency = "missing %s dependency in getter"

	if g.object == nil {
		return fmt.Errorf(missingDependency, "object storage")
	}

	return nil
}
