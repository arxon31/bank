package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DB DB `yaml:"db"`
}

type DB struct {
	URL string `env-required:"true" env:"PG_URL"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
