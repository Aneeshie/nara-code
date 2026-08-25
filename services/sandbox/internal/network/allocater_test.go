package network

import (
	"testing"
)

func TestAllocate(t *testing.T) {
	tests := []struct {
		name         string
		slot         int
		wantTap      string
		wantHostIP   string
		wantGuestIP  string
		wantGuestMAC string
	}{
		{
			name:         "Slot 0 - Start of range",
			slot:         0,
			wantTap:      "tap0",
			wantHostIP:   "172.16.0.1",
			wantGuestIP:  "172.16.0.2",
			wantGuestMAC: "06:00:ac:10:00:02",
		},
		{
			name:         "Slot 63 - Edge of first IP block",
			slot:         63,
			wantTap:      "tap63",
			wantHostIP:   "172.16.0.253",
			wantGuestIP:  "172.16.0.254",
			wantGuestMAC: "06:00:ac:10:00:fe",
		},
		{
			name:         "Slot 64 - Rolls over to next subnet block",
			slot:         64,
			wantTap:      "tap64",
			wantHostIP:   "172.16.1.1",
			wantGuestIP:  "172.16.1.2",
			wantGuestMAC: "06:00:ac:10:01:02",
		},
		{
			name:         "Slot 16383 - Max limit",
			slot:         16383,
			wantTap:      "tap16383",
			wantHostIP:   "172.16.255.253",
			wantGuestIP:  "172.16.255.254",
			wantGuestMAC: "06:00:ac:10:ff:fe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Allocate(tt.slot)
			if err != nil {
				t.Fatalf("Allocate(%d) returned unexpected error: %v", tt.slot, err)
			}

			if cfg.TapDevName != tt.wantTap {
				t.Errorf("TapDevName = %q; want %q", cfg.TapDevName, tt.wantTap)
			}
			if cfg.HostIP != tt.wantHostIP {
				t.Errorf("HostIP = %q; want %q", cfg.HostIP, tt.wantHostIP)
			}
			if cfg.GuestIP != tt.wantGuestIP {
				t.Errorf("GuestIP = %q; want %q", cfg.GuestIP, tt.wantGuestIP)
			}
			if cfg.GuestMAC != tt.wantGuestMAC {
				t.Errorf("GuestMAC = %q; want %q", cfg.GuestMAC, tt.wantGuestMAC)
			}
		})
	}
}

// TestAllocate_Errors checks that invalid slots correctly throw an error
func TestAllocate_Errors(t *testing.T) {
	tests := []struct {
		name string
		slot int
	}{
		{"Negative slot out of bounds", -1},
		{"Slot exactly at upper boundary limit", 16384},
		{"Slot way past upper boundary", 20000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Allocate(tt.slot)
			if err == nil {
				t.Errorf("Allocate(%d) expected an error for out-of-bounds slot, but got nil config: %+v", tt.slot, cfg)
			}
		})
	}
}

// TestIpToHex tests the internal unexported helper function
func TestIpToHex(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    string
		wantErr bool 
	}{
		{"Standard local IP", "192.168.1.1", "c0:a8:01:01", false},
		{"Your network IP block", "172.16.0.2", "ac:10:00:02", false},
		{"Invalid IP string", "999.999.999.999", "", true}, 
		{"Malformed text", "not-an-ip", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipToHex(tt.ip) 

			if (err != nil) != tt.wantErr {
				t.Fatalf("ipToHex(%q) unexpected error state: got error = %v, wantErr = %v", tt.ip, err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ipToHex(%q) = %q; want %q", tt.ip, got, tt.want)
			}
		})
	}
}
