package dependencies

import (
	"context"
	"time"

	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/cmd/mocks"
	"github.com/sanrinconr/storj-images/src/log"
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

	u, err := upload.NewAddImage(r, p.config.AllowedFormats, p.timer)
	if err != nil {
		panic(err)
	}

	return u
}

func (p packages) getAllLocations() mocks.GetterMock {
	return mocks.GetterMock(func(ctx context.Context) ([]domain.Location, error) {
		log.Info(ctx, "warning:using a mock to get images")

		return []domain.Location{
			{
				ID: "gatitohash",
				//nolint:lll // are a url of test.
				URL: "https://img.freepik.com/foto-gratis/gato-rojo-o-blanco-i-estudio-blanco_155003-13189.jpg?w=740&t=st=1668563236~exp=1668563836~hmac=fee01a12a48cc8abd9fb16a4e10a1be6ec7716c53474f92368d2200a38d81568",
			},
		}, nil
	})
}

func defaultTimer() func() time.Time {
	return time.Now
}
