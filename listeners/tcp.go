package listeners

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type TCPServer struct {
	address net.TCPAddr
}

var _ Listener = &TCPServer{}

func NewTCP(uri *url.URL) (*TCPServer, error) {
	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return nil, fmt.Errorf("could not parse port for TCP server (%s): %w", uri.String(), err)
	}

	address := net.TCPAddr{
		Port: port,
		IP:   net.ParseIP(uri.Hostname()),
	}

	return &TCPServer{
		address: address,
	}, nil
}

func (t *TCPServer) ListenAndServe(process ProcessMessage) error {
	server, err := net.ListenTCP("tcp", &t.address)
	if err != nil {
		return fmt.Errorf("cannot start TCP server (%s): %w", &t.address, err)
	}
	defer server.Close()

	p := make([]byte, maxUDPSize)
	for {
		connection, err := server.Accept()
		if err != nil {
			return fmt.Errorf("could not accept connection from TCP server (%s): %w", &t.address, err)
		}

		readLength, err := connection.Read(p[0:])
		if err != nil {
			return fmt.Errorf("could not read from TCP server (%s): %w", &t.address, err)
		}

		err = process(string(p[0:readLength]))
		if err != nil {
			return fmt.Errorf("could not process message from TCP server (%s): %w", &t.address, err)
		}
	}
}
