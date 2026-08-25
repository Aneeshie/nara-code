package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aneeshie/sandbox/internal/firecracker"
)

func main() {
	socketPath := "/tmp/firecracker.socket"

	//fail fast with a reasonable message if irecracker isnt up
	if _, err := os.Stat(socketPath); err != nil {
		log.Fatalf("no socket at %s — is firecracker running?", socketPath)
	}

	client := firecracker.NewClient(socketPath)

	err := client.SetLogger(&firecracker.LoggerConfig{
		LogPath:       "/home/aneeshie/sandbox/firecracker-go.log",
		Level:         firecracker.LevelFilterInfo,
		ShowLevel:     true,
		ShowLogOrigin: true,
	})
	if err != nil {
		log.Fatalf("SetLogger failed: %v", err)
	}

	fmt.Println("LOGGER OK")
}
