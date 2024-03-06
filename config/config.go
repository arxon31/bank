package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DB   DB   `yaml:"db"`
	AMQP AMQP `yaml:"amqp"`
}

type DB struct {
	URL string `env-required:"true" env:"PG_URL"`
}

type AMQP struct {
	URL       string `env-required:"true" env:"AMQP_URL"`
	QueueName string `env-required:"true" env:"AMQP_QUEUE_NAME"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
