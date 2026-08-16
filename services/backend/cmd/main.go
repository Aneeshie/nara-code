package main

import (
	"github.com/Aneeshie/nara-code/internal/config"
	"github.com/Aneeshie/nara-code/internal/grpc"
)

func main() {
	cfg := config.Load()

	grpc.CallAgent(cfg.AI_PORT)

}
