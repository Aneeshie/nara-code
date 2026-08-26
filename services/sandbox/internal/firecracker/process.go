package firecracker

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	cmd        *exec.Cmd
	socketPath string
}

func Start(socketPath string) (*Process, error) {
	_ = os.Remove(socketPath)

	cmd := exec.Command(
		"/home/aneeshie/sandbox/firecracker",
		"--api-sock",
		socketPath,
		"--enable-pci",
	)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start Firecracker: %w", err)
	}

	fmt.Println("firecracker started")

	process := &Process{
		cmd:        cmd,
		socketPath: socketPath,
	}

	if err := Poll(socketPath); err != nil {
		_ = process.Stop()
		return nil, err
	}

	return process, nil
}

func Poll(socketPath string) error {
	deadline := time.Now().Add(2 * time.Second)

	for {
		info, err := os.Stat(socketPath)

		if err == nil {
			if info.Mode()&os.ModeSocket != 0 {
				return nil
			}

			return fmt.Errorf("%s exists but is not a socket", socketPath)
		}

		if !os.IsNotExist(err) {
			return fmt.Errorf("stat Firecracker socket: %w", err)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for Firecracker socket")
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (p *Process) Wait() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	return p.cmd.Wait()
}

func (p *Process) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	if err := p.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("kill Firecracker: %w", err)
	}

	_, err := p.cmd.Process.Wait()

	_ = os.Remove(p.socketPath)

	return err
}
