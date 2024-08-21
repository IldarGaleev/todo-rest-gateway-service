// Package configapplication implements application configuration
package configapplication

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type AppConfig struct {
	EnvMode string `yaml:"env-mode" env-description:"" env:"ENV_MODE" env-default:"prod"`

	Grpc struct {
		Hostname string `yaml:"hostname" env-description:"" env:"HOSTNAME" env-default:"localhost"`
		Port     int    `yaml:"port" env-description:"" env:"PORT" env-default:"9090"`
	} `yaml:"grpc-client" env-prefix:"GRPC_"`

	Api struct {
		Hostname string `yaml:"hostname" env-description:"" env:"HOSTNAME" env-default:"localhost"`
		Port     int    `yaml:"port" env-description:"" env:"PORT" env-default:"8080"`
	} `yaml:"api" env-prefix:"API_"`
}

// MustLoadConfig Returns app configuration. Panic if failed
func MustLoadConfig(confPath string) *AppConfig {
	var appConf AppConfig

	err := cleanenv.ReadConfig(confPath, &appConf)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return &appConf
}
