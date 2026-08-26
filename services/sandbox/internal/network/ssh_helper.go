package network

import "fmt"

type SSHInfo struct {
	GuestIP string
	GuestMAC string 
	SSHClient string
}

func SSHForSlot(slot int) *SSHInfo {
	alloc, err := Allocate(slot)
	if err != nil || slot < 0 || slot >= maxSlots {
		return nil
	}

	return &SSHInfo{
        GuestIP:    alloc.GuestIP,
        GuestMAC:   alloc.GuestMAC,
        SSHClient:  fmt.Sprintf("ssh -i /home/aneeshie/sandbox/ubuntu-24.04.id_rsa root@%s", alloc.GuestIP),
    }
}


