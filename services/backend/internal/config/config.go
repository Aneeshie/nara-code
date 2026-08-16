package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT    string
	AI_PORT string
}

func NewConfig(PORT string, AI_PORT string) *Config {
	return &Config{
		PORT:    PORT,
		AI_PORT: AI_PORT,
	}
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env")
	}

	port := os.Getenv("PORT")

	aiAddr := os.Getenv("AI_PORT")

	cfg := NewConfig(port, aiAddr)

	return cfg
}
