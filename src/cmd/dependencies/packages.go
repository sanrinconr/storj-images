package dependencies

import (
	"time"

	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/getter"
	"github.com/sanrinconr/storj-images/src/upload"
)

type packages struct {
	config config.Config
	timer  func() time.Time
}

func newPackages(c config.Config, opts ...func() time.Time) packages {
	timer := defaultTimer()
	if len(opts) > 0 && opts[0] != nil {
		timer = opts[0]
	}

	return packages{c, timer}
}

func (p packages) uploadAddImage() upload.AddImage {
	r, err := upload.NewRepository(Mongo(p.config), Storj(p.config), defaultTimer())
	if err != nil {
		panic(err)
	}

	u, err := upload.NewAddImage(r, p.config.ImageAllowedFormats, p.timer)
	if err != nil {
		panic(err)
	}

	return u
}

func (p packages) getAllLocations() getter.Getter {
	g, err := getter.New(Mongo(p.config), Storj(p.config))
	if err != nil {
		panic(err)
	}

	return g
}

func defaultTimer() func() time.Time {
	return time.Now
}
