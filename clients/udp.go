package clients

import (
	"fmt"
	"net"
)

type UDPClient struct {
	connection net.Conn
}

func NewUDP(hostname string, port int) (*UDPClient, error) {
	connection, err :=  net.Dial("udp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		return nil, fmt.Errorf("could create UDP client (%s:%d): %w", hostname, port, err)
	}

	return &UDPClient{
		connection: connection,
	}, nil
}

func (u *UDPClient) WriteString(message string) error {
	length, err := fmt.Fprint(u.connection, message)
	if err != nil {
		return fmt.Errorf("could not write to UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	if length < len(message) {
		return fmt.Errorf("could not full message to UDP client (%s)", u.connection.LocalAddr())
	}

	return nil
}

func (u *UDPClient) Close() error {
	err := u.connection.Close()
	if err != nil {
		return fmt.Errorf("could not close UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	return nil
}
