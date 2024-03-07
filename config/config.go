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
	URL            string `env-required:"true" env:"AMQP_URL"`
	Exchange       string `env-required:"true" env:"AMQP_EXCHANGE"`
	Queue          string `env-required:"true" env:"AMQP_QUEUE"`
	RoutingKey     string `env-required:"true" env:"AMQP_ROUTING_KEY"`
	ConsumerTag    string `env-required:"true" env:"AMQP_CONSUMER_TAG"`
	WorkerPoolSize int    `env-required:"true" env:"AMQP_WORKER_POOL_SIZE"`
}

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
