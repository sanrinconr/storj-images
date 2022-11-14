package domain

import (
	"context"
)

// Uploader allow up the image into the cloud.
type Uploader interface {
	Upload(context.Context, Image) error
}

// Image are the representation of the input of a user.
// Have the raw image and the file name uploaded.
type Image struct {
	Raw  []byte
	Name string
}
