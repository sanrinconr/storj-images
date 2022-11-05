// Package repository give the more abstract methods to interact with the different storages.
package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/sanrinconr/storj-images/cmd/domain"
	"github.com/sanrinconr/storj-images/cmd/log"
	"github.com/sanrinconr/storj-images/cmd/repository/internal"
)

type (
	metadataInfra interface {
		Insert(context.Context, interface{}) error
	}

	objectInfra interface {
		Insert(context.Context, string, []byte) error
	}
)

// Image are a repository that can interact with a object storage (storj)
// and with a metadata storage (mongodb). Also include a timer to abstract
// the manage of the time in different time zones.
type Image struct {
	metadataInfra
	objectInfra
	timer func() time.Time
}

// NewImage create a new repository to abstract the interaction with multiple storages.
func NewImage(m metadataInfra, o objectInfra, t func() time.Time) (Image, error) {
	i := Image{
		metadataInfra: m,
		objectInfra:   o,
		timer:         t,
	}

	return i, i.validate()
}

// Insert add a new image, first adding into object storage and after adding in metadata storage.
func (i Image) Insert(ctx context.Context, img domain.Image) error {
	id := i.generateID()
	ext := img.Extension() // include point like: .xxx
	doc := internal.MetadataDoc{
		ID:               id,
		ObjectStorageKey: fmt.Sprintf("%s%s", id, ext),
		Name:             img.Name,
		CreatedAt:        i.timer(),
	}

	if err := i.metadataInfra.Insert(ctx, doc); err != nil {
		return err
	}

	log.Debug(ctx, fmt.Sprintf("image id '%s' added in metadata storage", doc.ID))

	if err := i.objectInfra.Insert(ctx, doc.ObjectStorageKey, img.Raw); err != nil {
		return err
	}

	log.Debug(ctx, fmt.Sprintf("image '%s' added in object storage", doc.ObjectStorageKey))

	return nil
}

// validate all dependencies of repository.
func (i Image) validate() error {
	const dependencyErr = "missing dependency %s in image repository"

	if i.metadataInfra == nil {
		return fmt.Errorf(dependencyErr, "metadata infra")
	}

	if i.objectInfra == nil {
		return fmt.Errorf(dependencyErr, "object infra")
	}

	if i.timer == nil {
		return fmt.Errorf(dependencyErr, "timer")
	}

	return nil
}

func (i Image) generateID() string {
	const len = 16 // https://en.wikipedia.org/wiki/Universally_unique_identifier

	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}
