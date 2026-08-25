package firecracker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	base       string
}

// ClientError is returned when Firecracker rejects a request with a non-2xx status.
type ClientError struct {
	Method       string
	Path         string
	StatusCode   int
	FaultMessage string
}

func (e *ClientError) Error() string {
	if e.FaultMessage != "" {
		return fmt.Sprintf("firecracker %s %s: %d: %s", e.Method, e.Path, e.StatusCode, e.FaultMessage)
	}
	return fmt.Sprintf("firecracker %s %s: unexpected status %d", e.Method, e.Path, e.StatusCode)
}

type LoggerConfig struct {
	// Named pipe or file used as output for logs.
	LogPath string `json:"log_path,omitempty"`
	// The level of the logger.
	Level LevelFilter `json:"level,omitempty"`
	// Whether to show the log level in the log.
	ShowLevel bool `json:"show_level,omitempty"`
	// Whether to show the log origin in the log.
	ShowLogOrigin bool `json:"show_log_origin,omitempty"`
}

type BootSourceConfig struct {
	// Path of the kernel image.
	KernelImagePath string `json:"kernel_image_path"`
	// Path of the initrd, if there is one.
	InitrdPath string `json:"initrd_path,omitempty"`
	// The boot arguments to pass to the kernel. If this field is uninitialized,
	// DEFAULT_KERNEL_CMDLINE is used.
	BootArgs string `json:"boot_args,omitempty"`
}

type RootFsConfig struct {
	// Unique identifier of the drive.
	DriveID string `json:"drive_id"`
	/// Part-UUID. Represents the unique id of the boot partition of this device. It is
	/// optional and it will be used only if the `is_root_device` field is true.
	PartitionUUID string `json:"partition_uuid,omitempty"`
	/// If set to true, it makes the current device the root block device.
	/// Setting this flag to true will mount the block device in the
	/// guest under /dev/vda unless the partuuid is present.
	IsRootDevice bool `json:"is_root_device"`
	// VirtioBlock specific fields
	/// If set to true, the drive is opened in read-only mode. Otherwise, the
	/// drive is opened as read-write.
	IsReadOnly bool `json:"is_read_only,omitempty"`
	/// Path of the drive.
	PathOnHost string `json:"path_on_host"`
	// VhostUserBlock specific fields
	/// Path to the vhost-user socket.
	SocketPath string `json:"socket_path,omitempty"`
}

type NetworkInterfaceConfig struct {
	/// ID of the guest network interface.
	IfaceID string `json:"iface_id"`
	/// Host level path for the guest network interface.
	HostDevName string `json:"host_dev_name"`
	/// Guest MAC address
	GuestMAC string `json:"guest_mac,omitempty"`
}

type LevelFilter string

const (
	LevelFilterOff   LevelFilter = "Off"
	LevelFilterTrace LevelFilter = "Trace"
	LevelFilterDebug LevelFilter = "Debug"
	LevelFilterInfo  LevelFilter = "Info"
	LevelFilterWarn  LevelFilter = "Warn"
	LevelFilterError LevelFilter = "Error"
)

type ActionType string

type ActionConfig struct {
	ActionType ActionType `json:"action_type"`
}

const (
	ActionFlushMetrics   ActionType = "FlushMetrics"
	ActionInstanceStart  ActionType = "InstanceStart"
	ActionSendCtrlAltDel ActionType = "SendCtrlAltDel"
)

func NewClient(socketPath string) *Client {
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   5 * time.Second,
	}

	return &Client{
		httpClient: &client,
		base:       "http://localhost/",
	}
}

// do sends one request to the Firecracker API socket. It returns a
// *ClientError for any non-2xx response, carrying the fault_message if present.
func (c *Client) do(method, path string, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal %s %s: %w", method, path, err)
	}

	req, err := http.NewRequest(method, c.base+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("build request %s %s: %w", method, path, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response %s %s: %w", method, path, err)
	}

	if resp.StatusCode >= 300 {
		var fault struct {
			FaultMessage string `json:"fault_message"`
		}
		_ = json.Unmarshal(respBody, &fault)

		return &ClientError{
			Method:       method,
			Path:         path,
			StatusCode:   resp.StatusCode,
			FaultMessage: fault.FaultMessage,
		}
	}

	return nil
}

func (c *Client) SetLogger(cfg LoggerConfig) error {
	return c.do(http.MethodPut, "logger", cfg)
}

func (c *Client) SetBootSource(cfg BootSourceConfig) error {
	return c.do(http.MethodPut, "boot-source", cfg)
}

func (c *Client) SetRootFs(cfg RootFsConfig) error {
	return c.do(http.MethodPut, "drives/rootfs", cfg)
}

func (c *Client) SetDrive(cfg RootFsConfig) error {
	return c.do(http.MethodPut, "drives/rootfs", cfg)
}

func (c *Client) SetNetworkInterface(cfg NetworkInterfaceConfig) error {
	return c.do(http.MethodPut, "network-interfaces/"+cfg.IfaceID, cfg)
}

func (c *Client) SetActions(cfg ActionConfig) error {
	return c.do(http.MethodPut, "actions", cfg)
}
