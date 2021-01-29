package listeners

import (
	"fmt"
	"net"
)

type ProcessMessage func(message string) error

type UDPServer struct {
	address net.UDPAddr
}

func NewUDP(hostname string, port int) (*UDPServer, error) {
	address := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(hostname),
	}

	return &UDPServer{
		address: address,
	}, nil
}

const maxUDPSize = 65507
func (u *UDPServer) ListenAndServe(process ProcessMessage) error {
	server, err := net.ListenUDP("udp", &u.address)
	if err != nil {
		return fmt.Errorf("cannot start UDP server (%s): %w", &u.address, err)
	}
	defer server.Close()

	p := make([]byte, maxUDPSize)
	for {
		readLength, _, err := server.ReadFromUDP(p[0:])
		if err != nil {
			return fmt.Errorf("could not read from UDP server (%s): %w", &u.address, err)
		}

		err = process(string(p[0:readLength]))
		if err != nil {
			return fmt.Errorf("could not process message from UDP server (%s): %w", &u.address, err)
		}
	}
}
