// Package core resides all the use cases of the application
package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/sanrinconr/storj-images/cmd/domain"
	"github.com/sanrinconr/storj-images/cmd/log"
)

type repository interface {
	Insert(context.Context, domain.Image) error
}

// AddImage has a repository to interact with storage and allowed formats to define what accept or not.
type AddImage struct {
	repository
	allowedFormats []string
}

// NewAddImage create a new use case for upload images.
func NewAddImage(r repository, f []string) (AddImage, error) {
	a := AddImage{
		repository:     r,
		allowedFormats: f,
	}

	return a, a.validate()
}

// Execute add a new image.
func (a AddImage) Execute(ctx context.Context, img domain.Image) error {
	if !a.allowedFormat(ctx, img) {
		return domain.InvalidFormatImageError(errors.New("invalid image format"))
	}

	if err := a.Insert(ctx, img); err != nil {
		return err
	}

	return nil
}

func (a AddImage) allowedFormat(ctx context.Context, i domain.Image) bool {
	received := i.Extension()
	for _, allowed := range a.allowedFormats {
		if received == allowed {
			return true
		}
	}

	log.Debug(ctx, fmt.Sprintf("image with format '%s' rejected", received))

	return false
}

func (a AddImage) validate() error {
	const dependencyErr = "missing %s in add image use case"

	if a.repository == nil {
		return fmt.Errorf(dependencyErr, "repository")
	}

	if len(a.allowedFormats) == 0 {
		return fmt.Errorf(dependencyErr, "allowed formats")
	}

	return nil
}
