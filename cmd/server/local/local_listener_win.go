//go:build windows

package local

import (
	"net"

	"github.com/Microsoft/go-winio"
)

func createLocalListener(path string) (net.Listener, error) {
	pipePath := `\\.\pipe\` + path

	return winio.ListenPipe(pipePath, nil)
}
