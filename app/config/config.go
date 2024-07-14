package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ServerHost        string `env:"SERVER_HOST"   env-default:"0.0.0.0"`
	ServerPort        int    `env:"SERVER_PORT"   env-default:"8080"`
	PostgresURI       string `env:"POSTGRES_URI"  env-required:""`
	Debug             bool   `env:"DEBUG"`
	DisableStacktrace bool   `env:"NO_STACKTRACE"`
}

func New() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
