package clients

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"net"
	"time"
)

type UDP struct {
	connection net.Conn
	config     *configuration
}

var _ Client = &UDP{}

func NewUDP(uri *url.URL) (*UDP, error) {
	connection, err := net.Dial("udp", fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port()))
	if err != nil {
		return nil, fmt.Errorf("could not create UDP client (%s:%s): %w", uri.Hostname(), uri.Port(), err)
	}

	config, err := newConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("could not validate configuration of UDP client (%s): %w", connection.LocalAddr(), err)
	}

	return &UDP{
		connection: connection,
		config: config,
	}, nil
}

func (u *UDP) Close() error {
	err := u.connection.Close()
	if err != nil {
		return fmt.Errorf("could not close UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	return nil
}

const maxUDPBufferSize = 65_507

func (u *UDP) ReadString() (string, error) {
	err := u.connection.SetReadDeadline(time.Now().Add(u.config.readDeadline))
	if err != nil {
		return "", fmt.Errorf("could not set write deadline to UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	buffer := make([]byte, maxUDPBufferSize)

	actualN, err := u.connection.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("could not read from UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	return string(buffer[:actualN]), nil
}

func (u *UDP) WriteString(message string) error {
	err := u.connection.SetWriteDeadline(time.Now().Add(u.config.writeDeadline))
	if err != nil {
		return fmt.Errorf("could not set write deadline to UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	length, err := u.connection.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("could not write to UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	if length < len(message) {
		return fmt.Errorf("could not full message to UDP client (%s)", u.connection.LocalAddr())
	}

	return nil
}
