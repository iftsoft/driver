package utils

import (
	"fmt"
	"net"
)

func isPortFree(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func getFreePort() (int, error) {
	// Listen on port 0 to let the OS choose a free port
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	// Get the port number assigned by the OS
	return l.Addr().(*net.TCPAddr).Port, nil
}
