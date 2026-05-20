//go:build !windows

package local

import (
	"net"
	"os"
)

func createLocalListener(path string) (net.Listener, error) {
	socketPath := path + ".sock"

	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			return nil, err
		}
	}

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	_ = os.Chmod(socketPath, 0600)

	return ln, nil
}
