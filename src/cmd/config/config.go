// Package config read configuration from yml files and environment variables.
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	// Prod define the production environment.
	Prod = "prod"
	// Develop define the development environment.
	Develop = "develop"
)

type (
	// Config general struct that have the configuration.
	Config struct {
		MongoMetadataDB     map[string]string `toml:"mongo_metadata_db"`
		StorjImagesDB       map[string]string `toml:"storj_images_db"`
		ImageAllowedFormats []string          `toml:"image_allowed_formats"`
	}
)

// ReadConfig from a file and return an object with all the configs.
func ReadConfig(env string) Config {
	location := fmt.Sprintf("./conf/%s.toml", ActualEnvironment())
	var c Config

	_, err := toml.DecodeFile(location, &c)
	if err != nil {
		panic(err)
	}

	return c
}

// ActualEnvironment return the scope running the api given by ENV variable.
func ActualEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" || env == Develop {
		return Develop
	}

	return Prod
}
