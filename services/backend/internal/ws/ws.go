package ws

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Aneeshie/nara-code/internal/grpc"
)

// func main() {
// 	log.SetFlags(0)

// 	err := run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// run starts a http.Server for the passed in address
// with all requests handled by echoServer.
func Run(agentService grpc.AgentClient) error {

	l, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		return err
	}
	log.Printf("listening on ws://%v", l.Addr())

	s := &http.Server{
		Handler: WebsocketServer{
			logf:         log.Printf,
			agentService: agentService,
		},
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
