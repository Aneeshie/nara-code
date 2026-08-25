package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Aneeshie/sandbox/internal/firecracker"
)

func main() {
	socketPath := "/tmp/firecracker.socket"

	//fail fast with a reasonable message if firecracker isnt up
	if _, err := os.Stat(socketPath); err != nil {
		log.Fatalf("no socket at %s — is firecracker running?", socketPath)
	}

	client := firecracker.NewClient(socketPath)

	err := client.SetLogger(firecracker.LoggerConfig{
		LogPath:       "/home/aneeshie/sandbox/firecracker-go.log",
		Level:         firecracker.LevelFilterInfo,
		ShowLevel:     true,
		ShowLogOrigin: true,
	})
	if err != nil {
		log.Fatalf("SetLogger failed: %v", err)
	}
	fmt.Println("LOGGER OK")

	err = client.SetBootSource(firecracker.BootSourceConfig{
		KernelImagePath: "/home/aneeshie/sandbox/vmlinux-6.18.41",
		BootArgs:        "console=ttyS0 reboot=k panic=1",
	})
	if err != nil {
		log.Fatalf("SetBootSource failed: %v", err)
	}
	fmt.Println("BOOT SOURCE OK")

	err = client.SetDrive(firecracker.RootFsConfig{
		DriveID:      "rootfs",
		PathOnHost:   "/home/aneeshie/sandbox/ubuntu-24.04.ext4",
		IsRootDevice: true,
		IsReadOnly:   false,
	})

	if err != nil {
		log.Fatalf("SetDrive failed: %v", err)
	}
	fmt.Println("DRIVE OK")

	err = client.SetNetworkInterface(firecracker.NetworkInterfaceConfig{
		IfaceID:     "net1",
		HostDevName: "tap0",
		GuestMAC:    "06:00:AC:10:00:02",
	})
	if err != nil {
		log.Fatalf("SetNetworkInterface failed: %v", err)
	}
	fmt.Println("NETWORK INTERFACE OK")

	err = client.SetActions(firecracker.ActionConfig{
		ActionType: firecracker.ActionInstanceStart,
	})
	if err != nil {
		log.Fatalf("SetActions failed: %v", err)
	}
	fmt.Println("ACTIONS OK")

}
