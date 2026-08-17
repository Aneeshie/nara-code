package ws

import (
	"io"
	"net/http"

	"github.com/Aneeshie/nara-code/internal/grpc"
	"github.com/Aneeshie/nara-code/internal/ws/dto"
	pb "github.com/Aneeshie/nara-code/proto"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WebsocketServer struct {
	// logf controls where logs are sent.
	logf         func(f string, v ...any)
	agentService grpc.AgentClient
}

func (s WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.CloseNow()

	for {
		var msg dto.Message

		err := wsjson.Read(r.Context(), c, &msg)
		if err != nil {
			return
		}
		msg.Role = dto.UserRole

		// call grpc functions
		stream, err := s.agentService.RunAgent(r.Context(), &msg)
		if err != nil {
			s.logf("error getting stream: %v", err)
			return
		}

		for {
			resp, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				s.logf("error getting stream: %v", err)
				return
			}

			respEvent := mapEvent(resp.GetEvent())

			// convert resp → WebSocket DTO
			agentMessage := dto.Message{
				Role:    dto.AssistantRole,
				Content: resp.GetMessage(),
				Event:   &respEvent,
			}
			// send to frontend
			err = wsjson.Write(r.Context(), c, agentMessage)

			if err != nil {
				s.logf("error sending stream to frontend: %v", err)
				return
			}

		}

	}
}

func mapEvent(event pb.EventType) dto.EventType {
	switch event {
	case pb.EventType_TEXT_DELTA:
		return dto.TextDeltaEvent

	case pb.EventType_TOOL_CALL:
		return dto.ToolEvent

	case pb.EventType_DONE:
		return dto.DoneEvent

	default:
		return ""
	}
}
