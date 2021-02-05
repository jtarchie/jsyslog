package listeners

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
)

type UDPServer struct {
	address net.UDPAddr
}

func NewUDP(uri *url.URL) (*UDPServer, error) {
	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return nil, fmt.Errorf("could not parse port for UDP server (%s): %w", uri.String(), err)
	}

	address := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(uri.Hostname()),
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
