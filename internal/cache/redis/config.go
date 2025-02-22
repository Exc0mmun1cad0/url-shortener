package redis

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// TODO: make it more secure (e. g. require password, protected mode and etc.)
// Configuration for redis connection.
type Config struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     int    `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password,omitempty" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

func MustLoad() Config {
	_ = godotenv.Load() // loading .env file

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Cannot read redis config from env vars: " + err.Error())
	}

	return cfg
}
