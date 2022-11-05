package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/domain"
	"github.com/sanrinconr/storj-images/cmd/log"
)

var errInternalFail = errors.New("retry later")

type executor interface {
	Execute(context.Context, domain.Image) error
}

// AddImage are the entrypoint that receive the data from the user.
type AddImage struct {
	executor
}

// NewAddImage create a controller to handle the user requests.
func NewAddImage(e executor) (AddImage, error) {
	adder := AddImage{
		executor: e,
	}

	return adder, adder.validate()
}

// AddImage are the entrypoint to save a new image into the app
// the controller can receive any file but.
func (a AddImage) AddImage(ctx *gin.Context) error {
	img, err := a.obtainImage(ctx)
	if err != nil {
		log.Error(ctx, err)

		if len(img.Raw) == 0 {
			return NewError(http.StatusBadRequest, errors.New("no param 'image' passed in the post form-data"))
		}

		return NewError(http.StatusInternalServerError, errInternalFail)
	}

	if len(img.Raw) == 0 {
		return NewError(http.StatusBadRequest, errors.New("no param 'image' is passed in post form-data"))
	}

	if err = a.executor.Execute(ctx, img); err != nil {
		var invalidFormatImage domain.InvalidFormatImageError
		if errors.As(err, &invalidFormatImage) {
			return NewError(http.StatusBadRequest, err)
		}

		log.Error(ctx, err)

		return NewError(http.StatusInternalServerError, errInternalFail)
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
	const dependencyErr = "missing  %s in add image controller"

	if a.executor == nil {
		return fmt.Errorf(dependencyErr, "executor")
	}

	return nil
}
