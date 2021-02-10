package servers

import (
	"bufio"
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"net"
)

type TCP struct {
	listener net.Listener
}

var _ Protocol = &TCP{}

func NewTCP(uri *url.URL) (*TCP, error) {
	listener, err := net.Listen("tcp", uri.Host)
	if err != nil {
		return nil, fmt.Errorf(
			"could not start listening on TCP (%s): %w",
			uri.Host,
			err,
		)
	}

	return &TCP{
		listener: listener,
	}, nil
}

func (t *TCP) Close() error {
	err := t.listener.Close()
	if err != nil {
		return fmt.Errorf(
			"could not close TCP (%s): %w",
			t.listener.Addr().String(),
			err,
		)
	}

	return nil
}

func (t *TCP) Listen() (Connection, error) {
	connection, err := t.listener.Accept()
	if err != nil {
		return nil, fmt.Errorf(
			"could not read from TCP (%s): %w",
			t.listener.Addr(),
			err,
		)
	}

	return &tcpConnection{
		connection: connection,
		buffer:     bufio.NewReader(connection),
	}, nil
}

func (t *TCP) LocalAddr() net.Addr {
	return t.listener.Addr()
}

type tcpConnection struct {
	buffer     *bufio.Reader
	connection net.Conn
}

func (t *tcpConnection) Write(p []byte) (int, error) {
	return t.connection.Write(p)
}

func (t *tcpConnection) Read(p []byte) (int, error) {
	return t.buffer.Read(p)
}

func (t *tcpConnection) LocalAddr() net.Addr {
	return t.connection.LocalAddr()
}

func (t *tcpConnection) RemoteAddr() net.Addr {
	return t.connection.RemoteAddr()
}
