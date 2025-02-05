package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"port"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

// Postgres represents connection to PostgreSQL databse
type Postgres struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name,omitempty"`
}

// MustLoad loads confiugration from .yaml file.
// If something bad occurs, it panics.
func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config with the following path does not exists: " + configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &config
}

// fetchConfigPath searches for config path.
// Firstly, it checks flag. If it's empty, it looks up in env var.
func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config-path", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
