package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		IsDebug    bool `env:"IS_DEBUG" env-default:"0"`
		HTTP       HTTP
		Log        Log
		Mongo      Mongo
		JWT        JWT
		WorkerPool WorkerPool
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"80"`
	}

	Log struct {
		Level  string `env:"LOG_LEVEL" env-default:"info"`
		Format string `env:"LOG_FORMAT" env-default:"json"`
	}

	Mongo struct {
		URL string `env:"MONGO_URL" env-required:"true"`
		DB  string `env:"MONGO_DB" env-required:"true"`
	}

	JWT struct {
		Secret string `env:"JWT_SECRET" env-required:"true"`
	}

	WorkerPool struct {
		Count int `env:"WORKER_POOL_COUNT" env-default:"10"`
	}
)

func New() (*Config, error) {
	const op = "config.New"

	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &cfg, nil
}
