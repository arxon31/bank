package cons

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	App  APP
	DB   DB
	AMQP AMQP
}

type APP struct {
	Env string `env-required:"true" env:"APP_ENV"`
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
