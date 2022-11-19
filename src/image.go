package domain

import (
	"context"
	"time"
)

// Uploader allow up the image into the cloud.
type Uploader interface {
	Upload(context.Context, Image) error
}

// Getter allow get images from the cloud.
type Getter interface {
	All(context.Context) ([]Location, error)
}

// Image are the representation of the input of a user.
// Have the raw image and the file name uploaded.
type Image struct {
	Raw  []byte
	Name string
}

// Location of image in internet .
type Location struct {
	ID        string
	Name      string
	URL       string
	CreatedAt time.Time
}
