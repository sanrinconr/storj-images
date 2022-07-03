package server

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type (
	// Config general struct that have the configuration.
	Config struct {
		Database
	}

	// Database related databases configuration.
	Database struct {
		IDS
		Photos
	}

	// IDS configuration related to the database of IDS.
	IDS struct {
		URL string `mapstructure:"url"`
	}

	// Photos configuration related to the storj photos.
	Photos struct {
		Project  string `mapstructure:"project"`
		Bucket   string `mapstructure:"bucket"`
		TokenENV string `mapstructure:"token_env"`
	}
)

func (c Config) validate() error {
	const MissingConfig = "missing config: %s"

	if c.Database.IDS.URL == "" {
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

// ReadConfig from a file and return an object with all the configs.
func ReadConfig() (Config, error) {
	var c Config

	v := viper.New()

	env := actualEnvironment()

	v.AddConfigPath("./conf/")
	v.SetConfigName(env)
	v.SetConfigType("yaml")

	log.Default().Printf("Using ENV: %s \n", env)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := v.Unmarshal(&c); err != nil {
		return Config{}, err
	}

	if err := c.validate(); err != nil {
		return Config{}, err
	}

	return c, nil
}

func actualEnvironment() string {
	const DefaultEnvironment = "develop"

	env := os.Getenv("ENV")
	if env == "" {
		return DefaultEnvironment
	}

	return env
}
