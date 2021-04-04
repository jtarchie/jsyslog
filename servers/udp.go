package servers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"io"
	"net"
)

type UDP struct {
	connection net.PacketConn
}

var _ Protocol = &UDP{}

func NewUDP(uri *url.URL) (*UDP, error) {
	connection, err := net.ListenPacket("udp", uri.Host)
	if err != nil {
		return nil, fmt.Errorf(
			"could not start listening on UDP (%s): %w",
			uri.Host,
			err,
		)
	}

	return &UDP{
		connection: connection,
	}, nil
}

func (u *UDP) Close() error {
	err := u.connection.Close()
	if err != nil {
		return fmt.Errorf(
			"could not close UDP (%s): %w",
			u.connection.LocalAddr().String(),
			err,
		)
	}

	return nil
}

func (u *UDP) Write(message []byte, connection Connection) error {
	_, err := u.connection.WriteTo(message, connection.RemoteAddr())
	if err != nil {
		return fmt.Errorf(
			"could not write to UDP (%s): %w",
			u.connection.LocalAddr(),
			err,
		)
	}

	return nil
}

func (u *UDP) Name() string {
	return "UDP"
}

const maxUDPBufferSize = 65_507

func (u *UDP) Listen() (Connection, error) {
	buffer := make([]byte, maxUDPBufferSize)

	actualN, addr, err := u.connection.ReadFrom(buffer[0:])
	if err != nil {
		return nil, fmt.Errorf(
			"could not read from UDP (%s): %w",
			u.connection.LocalAddr(),
			err,
		)
	}

	return &udpConnection{
		connection: u.connection,
		remote:     addr,
		buffer:     bufio.NewReader(bytes.NewReader(buffer[0:actualN])),
	}, nil
}

func (u *UDP) LocalAddr() net.Addr {
	return u.connection.LocalAddr()
}

type udpConnection struct {
	remote     net.Addr
	buffer     *bufio.Reader
	connection net.PacketConn
}

func (u *udpConnection) Close() error {
	return u.connection.Close()
}

func (u *udpConnection) Peek(i int) ([]byte, error) {
	return nil, io.EOF
}

func (u *udpConnection) Discard(i int) (int, error) {
	return 0, io.EOF
}

func (u *udpConnection) Write(p []byte) (int, error) {
	return u.connection.WriteTo(p, u.remote)
}

func (u *udpConnection) Read(p []byte) (int, error) {
	return u.buffer.Read(p)
}

func (u *udpConnection) LocalAddr() net.Addr {
	return u.connection.LocalAddr()
}

func (u *udpConnection) RemoteAddr() net.Addr {
	return u.remote
}
