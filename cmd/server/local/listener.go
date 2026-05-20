package local

import "net"

// CreateLocalListener creates a local listener for the given path.
// Compile for specific platform with flags
func CreateLocalListener(path string) (net.Listener, error) {
	return createLocalListener(path)
}
