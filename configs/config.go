package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Dsn    string
	Port   string
	ApiUrl string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Dsn:    os.Getenv("DSN"),
		Port:   os.Getenv("PORT"),
		ApiUrl: os.Getenv("API_URL"),
	}
}
