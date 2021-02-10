package clients

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"net"
	"time"
)

type UDP struct {
	connection net.Conn
}

var _ Client = &UDP{}

func NewUDP(uri *url.URL) (*UDP, error) {
	connection, err := net.Dial("udp", fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port()))
	if err != nil {
		return nil, fmt.Errorf("could not create UDP client (%s:%s): %w", uri.Hostname(), uri.Port(), err)
	}

	if timeout := uri.Query().Get("read-timeout"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse read timeout duration for UDP client (%s): %w",
				connection.LocalAddr().String(),
				err,
			)
		}

		err = connection.SetReadDeadline(time.Now().Add(duration))
		if err != nil {
			return nil, fmt.Errorf(
				"could not set read timeout duration for UDP client (%s): %w",
				connection.LocalAddr().String(),
				err,
			)
		}
	}

	if timeout := uri.Query().Get("write-timeout"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse write timeout duration for UDP client (%s): %w",
				connection.LocalAddr().String(),
				err,
			)
		}

		err = connection.SetWriteDeadline(time.Now().Add(duration))
		if err != nil {
			return nil, fmt.Errorf(
				"could not set write timeout duration for UDP client (%s): %w",
				connection.LocalAddr().String(),
				err,
			)
		}
	}

	return &UDP{
		connection: connection,
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
	buffer := make([]byte, maxUDPBufferSize)

	actualN, err := u.connection.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("could not read from UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	return string(buffer[:actualN]), nil
}

func (u *UDP) WriteString(message string) error {
	length, err := u.connection.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("could not write to UDP client (%s): %w", u.connection.LocalAddr(), err)
	}

	if length < len(message) {
		return fmt.Errorf("could not full message to UDP client (%s)", u.connection.LocalAddr())
	}

	return nil
}
