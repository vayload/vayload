//go:build windows

package client

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

func GetDialer(path string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return winio.DialPipeContext(ctx, `\\.\pipe\`+path)
	}
}
