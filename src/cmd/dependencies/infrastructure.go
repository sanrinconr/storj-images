package dependencies

import (
	"fmt"
	"os"

	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/storj"
)

// Storj resolve a storj infrastructure object.
func Storj(c config.Config) storj.Storj {
	t := os.Getenv(c.StorjImagesDB["token_env"])
	if t == "" {
		panic(fmt.Errorf("variable %s not is defined", c.StorjImagesDB["token_env"]))
	}

	s, err := storj.New(
		storj.WithAppAccess(t),
		storj.WithBucketName(c.StorjImagesDB["bucket"]),
		storj.WithProjectName(c.StorjImagesDB["project"]),
	)
	if err != nil {
		panic(err)
	}

	return s
}
