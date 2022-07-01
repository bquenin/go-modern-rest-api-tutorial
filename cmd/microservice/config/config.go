package config

import (
	"bytes"
	_ "embed"
	"strings"

	"github.com/spf13/viper"
)

//go:embed config.yml
var defaultConfiguration []byte

type Postgres struct {
	Host     string
	User     string
	Password string
}

type Config struct {
	Postgres *Postgres
}

func Read() (*Config, error) {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// Configuration file
	viper.SetConfigType("yml")

	// Read configuration
	if err := viper.ReadConfig(bytes.NewBuffer(defaultConfiguration)); err != nil {
		return nil, err
	}

	// Unmarshal the configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
