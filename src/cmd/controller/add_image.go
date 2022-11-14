// Package controller contains all the services exposed by the api with their respective logic (like errors management)
package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/log"
)

// AddImage are the entrypoint that receive the data from the user.
type AddImage struct {
	uploader domain.Uploader
}

// NewAddImage create a controller to handle the user requests.
func NewAddImage(u domain.Uploader) (func(*gin.Context) error, error) {
	a := AddImage{
		uploader: u,
	}
	if err := a.validate(); err != nil {
		return nil, err
	}

	return a.endpoint, nil
}

// endpoint are the entrypoint to save a new image into the app
// the controller can receive any file but.
func (a AddImage) endpoint(ctx *gin.Context) error {
	img, err := a.obtainImage(ctx)
	if err != nil {
		if len(img.Raw) == 0 {
			return NewError(http.StatusBadRequest, errors.New("no param 'image' passed in the post form-data"))
		}

		return NewError(http.StatusInternalServerError, err)
	}

	if len(img.Raw) == 0 {
		return NewError(http.StatusBadRequest, errors.New("no param 'image' is passed in post form-data"))
	}

	if err = a.uploader.Upload(ctx, img); err != nil {
		var invalidFormatImage domain.InvalidFormatImageError
		if errors.As(err, &invalidFormatImage) {
			return NewError(http.StatusBadRequest, err)
		}

		log.Error(ctx, err)

		return NewError(http.StatusInternalServerError, err)
	}

	return nil
}

func (a AddImage) obtainImage(ctx *gin.Context) (domain.Image, error) {
	file, header, err := ctx.Request.FormFile("image")
	defer func() {
		if file != nil {
			if err := file.Close(); err != nil {
				log.Info(ctx, err.Error())
			}
		}
	}()

	if err != nil {
		return domain.Image{}, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return domain.Image{}, err
	}

	return domain.Image{
		Raw:  buf.Bytes(),
		Name: header.Filename,
	}, nil
}

func (a AddImage) validate() error {
	const dependencyErr = "missing %s in add image controller"

	if a.uploader == nil {
		return fmt.Errorf(dependencyErr, "uploader")
	}

	return nil
}
