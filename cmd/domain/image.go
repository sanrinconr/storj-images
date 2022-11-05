package domain

import (
	"github.com/gabriel-vasile/mimetype"
)

// Image are the representation of the input of a user.
// Have the raw image and the file name uploaded.
type Image struct {
	Raw  []byte
	Name string
}

// Extension return the Extension of the file.
func (i Image) Extension() string {
	return mimetype.Detect(i.Raw).Extension()
}
