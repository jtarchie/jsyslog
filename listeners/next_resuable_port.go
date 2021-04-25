package listeners

import (
	"fmt"
	"net"
)

func NextReusablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("could not find the next available port: %w", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("could not test bind for listener: %w", err)
	}

	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}
