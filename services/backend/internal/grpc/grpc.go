package grpc

import (
	"context"

	"github.com/Aneeshie/nara-code/internal/ws/dto"
	pb "github.com/Aneeshie/nara-code/proto"
	"google.golang.org/grpc"
)

type AgentClient struct {
	client pb.AgentServiceClient
}

func NewAgentClient(conn *grpc.ClientConn) *AgentClient {
	return &AgentClient{
		client: pb.NewAgentServiceClient(conn),
	}
}

func (a *AgentClient) RunAgent(ctx context.Context, message *dto.Message) (grpc.ServerStreamingClient[pb.AgentResponse], error) {

	return a.client.RunAgent(
		ctx,
		&pb.AgentRequest{
			Message: message.Content,
		},
	)
}
