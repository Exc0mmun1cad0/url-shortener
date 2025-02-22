package postgres

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Configuration for postgres connection.
type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DBName   string `env:"POSTGRES_DB"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

func MustLoad() Config {
	_ = godotenv.Load() // loading .env file

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Cannot read postgres config from env vars: " + err.Error())
	}

	return cfg
}
