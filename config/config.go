package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string `env:"PORT" envDefault:":4000"`
	MongoURI string `env:"MONGO_URI,required"`
}

func New() Config {
	godotenv.Load()
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
