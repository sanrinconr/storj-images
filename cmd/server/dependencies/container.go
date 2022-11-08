// Package dependencies allow the inyection of dependencies in the respective
// parts of the every service.
package dependencies

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/cmd/controller"
)

// Container are the main dependency resolver, enter a request of resolve
// a controller and container provide all their dependencies.
type Container struct {
	infrastructure
	repository
	usecases
	Config
}

// New create a new container to resolve dependencies.
func New(env string) Container {
	conf := ReadConfig(env)
	i := newInfrastructure(ReadConfig(env))
	r := newRepository(&i)
	u := newUseCases(&r)

	return Container{
		infrastructure: i,
		repository:     r,
		usecases:       u,
		Config:         conf,
	}
}

// Ping are the basic service to validate if api is up or down.
func (c Container) Ping() func(*gin.Context) error {
	ctrl := controller.Ping{}

	return ctrl.Ping
}

// Error are a test endpoint to generate a controlled error.
func (c Container) Error() func(*gin.Context) error {
	ctrl := controller.ErrResponse{}

	return ctrl.Error
}

// AddImage generate the controller to add a image.
func (c Container) AddImage() func(*gin.Context) error {
	ctrl, err := controller.NewAddImage(
		c.usecases.AddImage(c.DefaultTimer(), c.config.AllowedFormats),
	)
	if err != nil {
		panic(err)
	}

	return ctrl.AddImage
}

// DefaultTimer generate the actual time when is called.
func (c Container) DefaultTimer() func() time.Time {
	return func() time.Time {
		return time.Now().UTC()
	}
}
