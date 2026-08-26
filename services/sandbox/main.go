package main

import (
	"log"

	"github.com/Aneeshie/sandbox/internal/firecracker"
	"github.com/Aneeshie/sandbox/internal/network"
)

func main() {
	socketPath := "/tmp/firecracker.socket"

	client := firecracker.NewClient(socketPath)

	alloc, err := network.Allocate(0)
	if err != nil {
		log.Fatal(err)
	}

	err = network.SetupTap(alloc)
	if err != nil {
		log.Fatal(err)
	}

	err = network.EnableInternetAccess(alloc, "enp3s0")
	if err != nil {
		log.Fatal(err)
	}

	err = network.EnableInternetAccess(alloc, "enp3s0")

	err = client.SetLogger(firecracker.LoggerConfig{
		LogPath: "/home/aneeshie/sandbox/firecracker-go.log",
		Level: "Debug",
		ShowLevel: true,
		ShowLogOrigin: true,
	})	

	if err != nil {
		log.Fatal(err)
	}

	err = client.SetBootSource(firecracker.BootSourceConfig{
		KernelImagePath: "/home/aneeshie/sandbox/vmlinux-6.18.41",
		BootArgs: "console=ttyS0 reboot=k panic=1",
	})
	
	if err != nil {
		log.Fatal(err)
	}

	err = client.SetRootFs(firecracker.RootFsConfig{
		DriveID: "rootfs",
		PathOnHost: "/home/aneeshie/sandbox/ubuntu-24.04.ext4",
		IsRootDevice: true,
		IsReadOnly: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = client.SetNetworkInterface(firecracker.NetworkInterfaceConfig{
		IfaceID: "net1",
		GuestMAC: alloc.GuestMAC,
		HostDevName: alloc.TapDevName,
	})

	if err != nil {
		log.Fatal(err)
	}

	client.Start()

}
