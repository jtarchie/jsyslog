package servers

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"io"
	"net"
)

type Connection interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	io.Reader
	io.Writer
}

type Protocol interface {
	Close() error
	Listen() (Connection, error)
	LocalAddr() net.Addr
}

type Handler interface {
	Receive(Connection) error
}

type Server struct {
	handler  Handler
	protocol Protocol
}

func NewServer(rawURL string, handler Handler) (*Server, error) {
	var protocol Protocol

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse server from url (%s): %w", rawURL, err)
	}

	switch uri.Scheme {
	case "udp":
		protocol, err = NewUDP(uri)
		if err != nil {
			return nil, fmt.Errorf("could not create UDP server: %w", err)
		}
	case "tcp":
		protocol, err = NewTCP(uri)
		if err != nil {
			return nil, fmt.Errorf("could not create TCP server: %w", err)
		}
	default:
		return nil, fmt.Errorf("could not create server from scheme %q", uri.Scheme)
	}

	return &Server{
		protocol: protocol,
		handler:  handler,
	}, nil
}

func (s *Server) ListenAndServe() error {
	protocol := s.protocol
	defer protocol.Close()

	for {
		connection, err := protocol.Listen()
		if err != nil {
			return fmt.Errorf("could not listen for connection (%s): %w", protocol.LocalAddr().String(), err)
		}

		err = s.handleConnection(connection)
		if err != nil {
			return err
		}
	}
}

func (s *Server) handleConnection(connection Connection) error {
	err := s.handler.Receive(connection)
	if err != nil {
		return fmt.Errorf(
			"could not handle packet from server: %w",
			err,
		)
	}

	return nil
}

func (s *Server) Close() error {
	return s.protocol.Close()
}
