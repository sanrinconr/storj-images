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

// Upload has a repository to interact with storage and allowed formats to define what accept or not.
type Upload struct {
	inserter       inserter
	allowedFormats []string
}

// NewAddImage create a new use case for upload images.
func NewAddImage(i inserter, f []string, timer func() time.Time) (Upload, error) {
	a := Upload{
		inserter:       i,
		allowedFormats: f,
	}

	return a, a.validate()
}

// Upload add a new image into the.
func (u Upload) Upload(ctx context.Context, img domain.Image) error {
	if !u.allowedFormat(ctx, img) {
		return domain.InvalidFormatImageError(errors.New("invalid image format"))
	}

	if err := u.inserter.Insert(ctx, img, u.extension(img)); err != nil {
		return err
	}

	return nil
}

func (u Upload) allowedFormat(ctx context.Context, i domain.Image) bool {
	received := u.extension(i)
	for _, allowed := range u.allowedFormats {
		if received == allowed {
			return true
		}
	}

	log.Debug(ctx, fmt.Sprintf("image with format '%s' rejected", received))

	return false
}

func (u Upload) extension(i domain.Image) string {
	return mimetype.Detect(i.Raw).Extension()
}

func (u Upload) validate() error {
	const dependencyErr = "missing %s in add image func"

	if u.inserter == nil {
		return fmt.Errorf(dependencyErr, "inserter")
	}

	if len(u.allowedFormats) == 0 {
		return fmt.Errorf(dependencyErr, "formats")
	}

	return nil
}
