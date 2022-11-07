package dependencies

import (
	"fmt"
	"os"

	infra "github.com/sanrinconr/storj-images/cmd/infrastructure"
)

// infrastructure needed infrastructure things.
// for what are pointers?, to make more easy the validation if exists or not an attribute.
type infrastructure struct {
	storj  *infra.Storj
	mongo  *infra.Mongo
	config Config
}

func newInfrastructure(c Config) infrastructure {
	return infrastructure{
		config: c,
	}
}

// Storj resolve a storj infrastructure object.
func (i *infrastructure) Storj() infra.Storj {
	if i.storj != nil {
		return *i.storj
	}

	t := os.Getenv(i.config.TokenENV)
	if t == "" {
		panic(fmt.Errorf("variable %s not is defined", i.config.TokenENV))
	}

	s, err := infra.NewStorj(
		infra.WithStorjAppAccess(t),
		infra.WithStorjBucketName(i.config.Bucket),
		infra.WithStorjProjectName(i.config.Project),
	)
	if err != nil {
		panic(err)
	}

	i.storj = &s

	return *i.storj
}

// Storj resolve a storj infrastructure object.
func (i *infrastructure) Mongo() infra.Mongo {
	if i.mongo != nil {
		return *i.mongo
	}

	mongo, err := infra.NewMongo(
		os.Getenv(i.config.IDS.URLENV),
		i.config.IDS.Database,
		i.config.IDS.Collection,
		os.Getenv(i.config.IDS.UserENV),
		os.Getenv(i.config.IDS.PasswordENV),
	)
	if err != nil {
		panic(err)
	}

	i.mongo = &mongo

	return mongo
}
