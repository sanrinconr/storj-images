// Package dependencies is responsible of resolve and inject dependencies.
package dependencies

import (
	"github.com/gin-gonic/gin"
	"github.com/sanrinconr/storj-images/src/cmd/config"
	"github.com/sanrinconr/storj-images/src/cmd/controller"
)

// Resolver generate all controllers.
type Resolver struct {
	packages packages
}

// New create a resolver that is capable of create controllers.
func New() Resolver {
	c := config.ReadConfig(config.ActualEnvironment())
	p := newPackages(c)

	return Resolver{
		packages: p,
	}
}

// Ping check health status of the api.
func (r Resolver) Ping() func(*gin.Context) error {
	return controller.NewPing()
}

// UploadImage is controller that allow upload images in the cloud.
func (r Resolver) UploadImage() func(*gin.Context) error {
	ctrl, err := controller.NewAddImage(r.packages.uploadAddImage())
	if err != nil {
		panic(err)
	}

	return ctrl
}

// GetAllLocations resolve the controller to get images from the cloud.
func (r Resolver) GetAllLocations() func(*gin.Context) error {
	ctrl, err := controller.NewGetAllLocations(r.packages.getAllLocations())
	if err != nil {
		panic(err)
	}

	return ctrl
}
