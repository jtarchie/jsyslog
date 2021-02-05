package clients

import (
	"fmt"
	"net"
	"net/url"
)

type UDPClient struct {
	connection net.Conn
}

var _ Client = &UDPClient{}

func NewUDP(uri *url.URL) (*UDPClient, error) {
	connection, err := net.Dial("udp", fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port()))
	if err != nil {
		return nil, fmt.Errorf("could not create UDP client (%s:%s): %w", uri.Hostname(), uri.Port(), err)
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
