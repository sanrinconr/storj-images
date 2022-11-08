package dependencies

import (
	"time"

	repo "github.com/sanrinconr/storj-images/cmd/repository"
)

type repository struct {
	insertImage    *repo.Image
	infrastructure *infrastructure
}

func newRepository(i *infrastructure) repository {
	return repository{
		infrastructure: i,
	}
}

func (r *repository) InsertImage(t func() time.Time) repo.Image {
	if r.insertImage != nil {
		// pointer is use only to validate if repo are already defined or not
		return *r.insertImage
	}

	res, err := repo.NewImage(r.infrastructure.Mongo(), r.infrastructure.Storj(), t)
	if err != nil {
		panic(err)
	}

	r.insertImage = &res

	return *r.insertImage
}
