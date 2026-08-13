package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func NewConfig(PORT string) *Config {
	return &Config{
		PORT: PORT,
	}
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env")
	}

	port := os.Getenv("PORT")

	cfg := NewConfig(port)

	return cfg
}
