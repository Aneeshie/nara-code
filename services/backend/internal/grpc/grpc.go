package grpc

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/Aneeshie/nara-code/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CallAgent(aiAddr string) {
	conn, err := grpc.NewClient("localhost:"+aiAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}

	defer conn.Close()

	//get the typed client
	client := pb.NewAgentServiceClient(conn)

	//run the agent
	stream, err := client.RunAgent(context.Background(), &pb.AgentRequest{Message: "hi"})
	if err != nil {
		log.Fatalf("RunAgent failed: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("recv error: %v", err)
		}

		fmt.Println(resp.GetMessage()) // print each streamed chunk
	}
}
