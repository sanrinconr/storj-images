// Package mocks have the test implementations for domain interfaces.
package mocks

import (
	"context"
	"errors"

	domain "github.com/sanrinconr/storj-images/src"
)

// GetterMock simulate the response of location of images in the cloud.
type GetterMock func(context.Context) ([]domain.Location, error)

// GetAll can simulate any behaviour.
func (g GetterMock) GetAll(ctx context.Context) ([]domain.Location, error) {
	if g == nil {
		return nil, errors.New("default mock error")
	}

	return g(ctx)
}