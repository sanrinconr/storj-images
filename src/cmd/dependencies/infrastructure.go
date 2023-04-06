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
	t := os.Getenv(c.TokenENV)
	if t == "" {
		panic(fmt.Errorf("variable %s not is defined", c.TokenENV))
	}

	s, err := storj.New(
		storj.WithAppAccess(t),
		storj.WithBucketName(c.Bucket),
		storj.WithProjectName(c.Project),
	)
	if err != nil {
		panic(err)
	}

	return s
}

// Mongo resolve a mongo infrastructure object to save documents.
func Mongo(c config.Config) mongo.Mongo {
	m, err := mongo.NewMongo(
		os.Getenv(c.IDS.URLENV),
		c.IDS.Database,
		c.IDS.Collection,
		os.Getenv(c.IDS.UserENV),
		os.Getenv(c.IDS.PasswordENV),
	)
	if err != nil {
		panic(err)
	}

	return m
}
