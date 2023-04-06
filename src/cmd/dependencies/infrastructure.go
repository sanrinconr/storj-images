package dependencies

import (
	"fmt"
	"os"

	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/mongo"
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

// Mongo resolve a mongo infrastructure object to save documents.
func Mongo(c config.Config) mongo.Mongo {
	m, err := mongo.NewMongo(
		os.Getenv(c.MongoMetadataDB["url_env"]),
		c.MongoMetadataDB["database"],
		c.MongoMetadataDB["collection"],
		os.Getenv(c.MongoMetadataDB["user_env"]),
		os.Getenv(c.MongoMetadataDB["password_env"]),
	)
	if err != nil {
		panic(err)
	}

	return m
}
