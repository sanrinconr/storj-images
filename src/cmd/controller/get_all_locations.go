package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/sanrinconr/storj-images/src"
	"github.com/sanrinconr/storj-images/src/cmd/controller/internal"
)

// GetAllLocations return all the images saved in the cloud.
type GetAllLocations struct {
	getter domain.Getter
}

// NewGetAllLocations create the controller to dowload images from the cloud.
func NewGetAllLocations(getter domain.Getter) (func(*gin.Context) error, error) {
	g := GetAllLocations{
		getter: getter,
	}

	return g.endpoint, g.validate()
}

func (g GetAllLocations) endpoint(ctx *gin.Context) error {
	locations, err := g.getter.GetAll(ctx)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}

	response := g.domainToResponse(locations)

	ctx.JSON(http.StatusOK, response)

	return nil
}

func (g GetAllLocations) domainToResponse(l []domain.Location) []internal.Location {
	r := make([]internal.Location, len(l))

	for i := range l {
		r[i] = internal.Location{
			ID:  l[i].ID,
			URL: l[i].URL,
		}
	}

	return r
}

func (g GetAllLocations) validate() error {
	const missingDependency = "missing %s in controller get all images"

	if g.getter == nil {
		return fmt.Errorf(missingDependency, "getter")
	}

	return nil
}
