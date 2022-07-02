package server

import (
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
	}

	// IDS configuration related to the database of IDS.
	IDS struct {
		URL string `mapstructure:"url"`
	}
)

func readConfig() (Config, error) {
	var c Config

	v := viper.New()

	v.AddConfigPath("./conf/")
	v.SetConfigName(actualEnvironment())
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := v.Unmarshal(&c); err != nil {
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
