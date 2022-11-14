// Package upload give the more abstract methods to interact with the different storages.
package upload

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/log"
	"github.com/sanrinconr/storj-images/src/upload/internal"
)

type (
	// MetadataInfra are the abstraction of the storage used to save metadata like a document.
	MetadataInfra interface {
		Insert(context.Context, interface{}) error
	}

	// ObjectInfra are the abstraction of the storage used to save a object like an image.
	ObjectInfra interface {
		Insert(context.Context, string, []byte) error
	}
)

// Repository can interact with a object storage (storj)
// and with a metadata storage (mongodb). Also include a timer to abstract
// the manage of the time in different time zones.
type Repository struct {
	MetadataInfra
	ObjectInfra
	timer func() time.Time
}

// NewRepository create a new repository to abstract the interaction with multiple storages.
func NewRepository(m MetadataInfra, o ObjectInfra, t func() time.Time) (Repository, error) {
	r := Repository{
		MetadataInfra: m,
		ObjectInfra:   o,
		timer:         t,
	}

	return r, r.validate()
}

// Insert add a new image, first adding into object storage and after adding in metadata storage.
func (r Repository) Insert(ctx context.Context, img domain.Image, ext string) error {
	id := r.generateID()
	doc := internal.MetadataDoc{
		ID:               id,
		ObjectStorageKey: fmt.Sprintf("%s%s", id, ext),
		Name:             img.Name,
		CreatedAt:        r.timer(),
	}

	if err := r.MetadataInfra.Insert(ctx, doc); err != nil {
		return err
	}

	log.Info(ctx, fmt.Sprintf("image id '%s' added in metadata storage", doc.ID))

	if err := r.ObjectInfra.Insert(ctx, doc.ObjectStorageKey, img.Raw); err != nil {
		return err
	}

	log.Info(ctx, fmt.Sprintf("image '%s' added in object storage", doc.ObjectStorageKey))

	return nil
}

// validate all dependencies of repository.
func (r Repository) validate() error {
	const dependencyErr = "missing dependency %s in image repository"

	if r.MetadataInfra == nil {
		return fmt.Errorf(dependencyErr, "metadata infra")
	}

	if r.ObjectInfra == nil {
		return fmt.Errorf(dependencyErr, "object infra")
	}

	if r.timer == nil {
		return fmt.Errorf(dependencyErr, "timer")
	}

	return nil
}

func (r Repository) generateID() string {
	const len = 16 // https://en.wikipedia.org/wiki/Universally_unique_identifier

	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}
