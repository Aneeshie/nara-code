package main

import (
	"log"

	"github.com/Aneeshie/nara-code/internal/config"
	"github.com/Aneeshie/nara-code/internal/ws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	g "github.com/Aneeshie/nara-code/internal/grpc"
)

func main() {
	cfg := config.Load()

	client, err := grpc.NewClient("localhost:"+cfg.GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error connecting grpc server: %v", err)
	}

	agentClient := g.NewAgentClient(client)

	ws.Run(*agentClient)

}
