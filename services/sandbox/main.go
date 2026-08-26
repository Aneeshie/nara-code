package main

import (
	"fmt"
	"log"

	"github.com/Aneeshie/sandbox/internal/firecracker"
	"github.com/Aneeshie/sandbox/internal/network"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}

}	

func run() error {
	socketPath := "/tmp/firecracker.socket"

	client := firecracker.NewClient(socketPath)

	slot := network.NextAvailable()

	alloc, err := network.Allocate(slot)
	if err != nil {
		return err
	}

	defer network.Release(slot)

	err = network.SetupTap(alloc)
	if err != nil {
		return err
	}

	err = network.EnableInternetAccess(alloc, "enp3s0")
	if err != nil {
		return err
	}

	proc, err := firecracker.Start(socketPath)
	if err != nil {
		return err
	}

	defer proc.Stop()


	err = client.SetLogger(firecracker.LoggerConfig{
		LogPath: "/home/aneeshie/sandbox/firecracker-go.log",
		Level: firecracker.LevelFilterDebug,
		ShowLevel: true,
		ShowLogOrigin: true,
	})	

	if err != nil {
		proc.Stop()
		return err
	}

	err = client.SetBootSource(firecracker.BootSourceConfig{
		KernelImagePath: "/home/aneeshie/sandbox/vmlinux-6.18.41",
		BootArgs: "console=ttyS0 reboot=k panic=1",
	})

	if err != nil {
		proc.Stop()
		return err
	}

	err = client.SetRootFs(firecracker.RootFsConfig{
		DriveID: "rootfs",
		PathOnHost: "/home/aneeshie/sandbox/ubuntu-24.04.ext4",
		IsRootDevice: true,
		IsReadOnly: false,
	})

	if err != nil {
		proc.Stop()
		return err
	}

	err = client.SetNetworkInterface(firecracker.NetworkInterfaceConfig{
		IfaceID: "net1",
		GuestMAC: alloc.GuestMAC,
		HostDevName: alloc.TapDevName,
	})

	if err != nil {
		proc.Stop()
		return err
	}

	err = client.Start()
	if err != nil {
		proc.Stop()
		return err
	}

	fmt.Println("VM STARTED")
	return nil
}
