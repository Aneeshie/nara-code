package network

import (
	"fmt"
	"net"
	"strconv"

)

type MVMNetworkConfig struct{
	TapDevName string
	GuestIP string
	HostIP string
	GuestMAC string
	Mask string
}

func (c *MVMNetworkConfig) TapAddrWithMask() string {
	return c.HostIP + "/30"
}


func Allocate(slot int) (*MVMNetworkConfig, error) {
	if slot < 0 || slot >= 16384 { 
		return nil, fmt.Errorf("the slot range should be between 0 and 16383")
	}

	tapDevName := "tap" + strconv.Itoa(slot)

	host_ip := fmt.Sprintf("172.16.%d.%d", slot/64, (4*slot+1)%256)
	guest_ip := fmt.Sprintf("172.16.%d.%d", slot/64, (4*slot+2)%256)

	guest_hex_ip, err := ipToHex(guest_ip)
	if err != nil {
		return nil, fmt.Errorf("Invalid Generated Guest IP address: %v", err)
	}
	
	guestMac := "06:00:" + guest_hex_ip

	return &MVMNetworkConfig{
		TapDevName: tapDevName,
		GuestIP:    guest_ip,
		HostIP:     host_ip,
		GuestMAC:   guestMac,
		Mask: "/30",
	}, nil
}

func ipToHex(ipStr string) (string, error) { ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", fmt.Errorf("Invalid IP")
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return "", fmt.Errorf("Not a valid IPv4 address")
	}

	hexStr := fmt.Sprintf("%02x:%02x:%02x:%02x", ipv4[0], ipv4[1], ipv4[2], ipv4[3])
	return hexStr, nil
}
