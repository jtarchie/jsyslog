package clients

import (
	"fmt"
	"net"
	"net/url"
)

type TCPClient struct {
	connection net.Conn
}
var  _ Client = &TCPClient{}

func NewTCP(uri *url.URL) (*TCPClient, error) {
	connection, err :=  net.Dial("tcp", fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port()))
	if err != nil {
		return nil, fmt.Errorf("could not create TCP client (%s:%s): %w", uri.Hostname(), uri.Port(), err)
	}

	return &TCPClient{
		connection: connection,
	}, nil
}

func (t *TCPClient) WriteString(message string) error {
	length, err := fmt.Fprint(t.connection, message)
	if err != nil {
		return fmt.Errorf("could not write to TCP client (%s): %w", t.connection.LocalAddr(), err)
	}

	if length < len(message) {
		return fmt.Errorf("could not full message to TCP client (%s)", t.connection.LocalAddr())
	}

	return nil
}

func (t *TCPClient) Close() error {
	err := t.connection.Close()
	if err != nil {
		return fmt.Errorf("could not close TCP client (%s): %w", t.connection.LocalAddr(), err)
	}

	return nil
}
