package dependencies

import (
	"time"

	"github.com/sanrinconr/storj-images/cmd/core"
)

type usecases struct {
	addImage   *core.AddImage
	repository *repository
}

func newUseCases(r *repository) usecases {
	return usecases{
		repository: r,
	}
}

func (u *usecases) AddImage(t func() time.Time, allowedFormats []string) core.AddImage {
	if u.addImage != nil {
		return *u.addImage
	}

	c, err := core.NewAddImage(u.repository.InsertImage(t), allowedFormats)
	if err != nil {
		panic(err)
	}

	u.addImage = &c

	return c
}
