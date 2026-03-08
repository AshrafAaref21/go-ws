package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `env:"HTTP_ADDRESS" envDefault:"localhost:8082"`
}

type Config struct {
	Env    string `env:"ENV" envDefault:"dev"`
	DBPath string `env:"DB_PATH" envDefault:"sqlite/dev"`
	DBName string `env:"DB_NAME" envDefault:"api.db"`
	HTTPServer
	JWTKey string `env:"JWT_KEY" envDefault:"supersecretkey"`
}

func LoadConfig() *Config {
	var cfg Config
	var envPath string

	flag.StringVar(&envPath, "config", "", "path to .env file")
	flag.Parse()

	if envPath == "" {
		envPath = os.Getenv("CONFIG_PATH")
	}

	if envPath == "" {
		envPath = "config/dev.env"
	}

	err := cleanenv.ReadConfig(envPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
