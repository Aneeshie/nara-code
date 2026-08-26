package network

import (
	"fmt"
	"os/exec"
)

func SetupTap(cfg *MVMNetworkConfig) error {
	// 1. ip link del <TapDevName>
	exec.Command("ip", "link", "del", cfg.TapDevName).Run()
	// 2. ip tuntap add <TapDevName> mode tap
	if err := exec.Command("ip", "tuntap", "add", cfg.TapDevName, "mode", "tap").Run(); err != nil {
		return err
	}
	// 3. ip addr add <cfg.TapAddrWithMask()> dev <TapDevName>
	if err := exec.Command("ip", "addr", "add", cfg.TapAddrWithMask(), "dev", cfg.TapDevName).Run(); err != nil {
		return err
	}
	// 4. ip link set <TapDevName> up
	if err := exec.Command("ip", "link", "set", cfg.TapDevName, "up").Run(); err != nil {
		return err
	}
	
	return nil	
}

func EnableInternetAccess(cfg *MVMNetworkConfig, hostIface string) error {
	// 1. Enable IPv4 forwarding
	if err := exec.Command(
		"sysctl",
		"-w",
		"net.ipv4.ip_forward=1",
	).Run(); err != nil {
		return fmt.Errorf("enable IP forwarding: %w", err)
	}

	// 2. Create nftables table
	// Ignore error because the table may already exist.
	exec.Command(
		"nft",
		"add",
		"table",
		"ip",
		"firecracker",
	).Run()

	// 3. Create postrouting NAT chain
	// Ignore error because it may already exist.
	exec.Command(
		"nft",
		"add",
		"chain",
		"ip",
		"firecracker",
		"postrouting",
		"{",
		"type",
		"nat",
		"hook",
		"postrouting",
		"priority",
		"100",
		";",
		"}",
	).Run()

	// Create forward filter chain
	// Ignore error because it may already exist.
	exec.Command(
		"nft",
		"add",
		"chain",
		"ip",
		"firecracker",
		"forward",
		"{",
		"type",
		"filter",
		"hook",
		"forward",
		"priority",
		"0",
		";",
		"}",
	).Run()

	// 4. Masquerade VM traffic leaving through the host interface.
	if err := exec.Command(
		"nft",
		"add",
		"rule",
		"ip",
		"firecracker",
		"postrouting",
		"ip",
		"saddr",
		cfg.GuestIP,
		"oifname",
		hostIface,
		"masquerade",
	).Run(); err != nil {
		return fmt.Errorf("add masquerade rule: %w", err)
	}

	// 5. Allow traffic from the VM's TAP interface
	// to leave through the host's Internet interface.
	if err := exec.Command(
		"nft",
		"add",
		"rule",
		"ip",
		"firecracker",
		"forward",
		"iifname",
		cfg.TapDevName,
		"oifname",
		hostIface,
		"accept",
	).Run(); err != nil {
		return fmt.Errorf("add forward rule: %w", err)
	}

	return nil
}

func SetupFirecrackerProcess() {

}
