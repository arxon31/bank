package prod

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env  string `env-required:"true" env:"PRODUCER_ENV"`
	AMQP AMQP
}

type AMQP struct {
	URL        string `env-required:"true" env:"AMQP_URL"`
	Exchange   string `env-required:"true" env:"AMQP_EXCHANGE"`
	RoutingKey string `env-required:"true" env:"AMQP_ROUTING_KEY"`
}

func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
