//go:build !windows

package client

import (
	"context"
	"net"
)

func GetDialer(path string) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", path+".sock")
	}
}
