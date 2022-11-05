// Package internal are the representations of the user inputs
package internal

import "mime/multipart"

// FormImage are the file received when a image is loaded.
type FormImage struct {
	Image *multipart.FileHeader `form:"image" binding:"required"`
}
