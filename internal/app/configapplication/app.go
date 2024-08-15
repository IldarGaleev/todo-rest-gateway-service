// Package configapplication implements application configuration
package configapplication

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	EnvMode      string `env:"ENV_MODE" env-default:"prod"`
	GrpcHostname string `env:"GRPC_HOSTNAME" env-default:"localhost"`
	GrpcPort     int    `env:"GRPC_PORT" env-default:"9090"`

	APIHostname string `env:"API_HOSTNAME" env-default:"localhost"`
	APIPort     int    `env:"GRPC_PORT" env-default:"8080"`
}

// MustLoadConfig Returns app configuration. Panic if failed
func MustLoadConfig() *AppConfig {
	//TODO: add loading config from file
	var appConf AppConfig
	err := cleanenv.ReadEnv(&appConf)
	if err != nil {
		panic(err)
	}
	return &appConf
}
