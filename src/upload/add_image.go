// Package upload resides all the use cases of the application
package upload

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gabriel-vasile/mimetype"
	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/log"
)

type inserter interface {
	Insert(context.Context, domain.Image, string) error
}

// AddImage has a repository to interact with storage and allowed formats to define what accept or not.
type AddImage struct {
	inserter       inserter
	allowedFormats []string
}

// NewAddImage create a new use case for upload images.
func NewAddImage(i inserter, f []string, timer func() time.Time) (AddImage, error) {
	a := AddImage{
		inserter:       i,
		allowedFormats: f,
	}

	return a, a.validate()
}

// Upload add a new image into the.
func (a AddImage) Upload(ctx context.Context, img domain.Image) error {
	if !a.allowedFormat(ctx, img) {
		return domain.InvalidFormatImageError(errors.New("invalid image format"))
	}

	if err := a.inserter.Insert(ctx, img, a.extension(img)); err != nil {
		return err
	}

	return nil
}

func (a AddImage) allowedFormat(ctx context.Context, i domain.Image) bool {
	received := a.extension(i)
	for _, allowed := range a.allowedFormats {
		if received == allowed {
			return true
		}
	}

	log.Debug(ctx, fmt.Sprintf("image with format '%s' rejected", received))

	return false
}

func (a AddImage) extension(i domain.Image) string {
	return mimetype.Detect(i.Raw).Extension()
}

func (a AddImage) validate() error {
	const dependencyErr = "missing %s in add image func"

	if a.inserter == nil {
		return fmt.Errorf(dependencyErr, "inserter")
	}

	if len(a.allowedFormats) == 0 {
		return fmt.Errorf(dependencyErr, "formats")
	}

	return nil
}
