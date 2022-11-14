// Package config read configuration from yml files and environment variables.
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	prod    = "prod"
	develop = "develop"
)

type (
	// Config general struct that have the configuration.
	Config struct {
		Database
		Image
	}

	// Database related databases configuration.
	Database struct {
		IDS
		Photos
	}

	// IDS configuration related to the database of IDS.
	IDS struct {
		URLENV      string `mapstructure:"url_env"`
		Database    string `mapstructure:"database"`
		Collection  string `mapstructure:"collection"`
		UserENV     string `mapstructure:"user_env"`
		PasswordENV string `mapstructure:"password_env"`
	}

	// Photos configuration related to the storj photos.
	Photos struct {
		Project  string `mapstructure:"project"`
		Bucket   string `mapstructure:"bucket"`
		TokenENV string `mapstructure:"token_env"`
	}

	// Image allow manage the available images format to be uploaded.
	Image struct {
		AllowedFormats []string `mapstructure:"allowed_formats"`
	}
)


// ReadConfig from a file and return an object with all the configs.
func ReadConfig(env string) Config {
	var c Config

	v := viper.New()

	v.AddConfigPath("./conf/")
	v.SetConfigName(env)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&c); err != nil {
		panic(err)
	}

	if err := c.validate(); err != nil {
		panic(err)
	}

	return c
}

// ActualEnvironment return the scope running the api given by ENV variable.
func ActualEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" || env == develop {
		return develop
	}

	return prod
}

func (c Config) validate() error {
	const MissingConfig = "missing config: %s"

	if c.Database.IDS.URLENV == "" {
		return fmt.Errorf(MissingConfig, "Database.IDS.URL")
	}

	if c.Database.Photos.Project == "" {
		return fmt.Errorf(MissingConfig, "Database.Photos.Project")
	}

	if c.Database.Photos.Bucket == "" {
		return fmt.Errorf(MissingConfig, "Database.Photos.Bucket")
	}

	if c.Database.Photos.TokenENV == "" {
		return fmt.Errorf(MissingConfig, "Database.Photos.TokenENV")
	}

	return nil
}
