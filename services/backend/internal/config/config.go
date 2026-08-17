package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WS_PORT   string
	GRPC_PORT string
}

func NewConfig(WS_PORT string, GRPC_PORT string) *Config {
	return &Config{
		WS_PORT:   WS_PORT,
		GRPC_PORT: GRPC_PORT,
	}
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env")
	}

	wsPort := os.Getenv("WS_PORT")

	gRPCPORT := os.Getenv("GRPC_PORT")

	cfg := NewConfig(wsPort, gRPCPORT)

	return cfg
}
