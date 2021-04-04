package servers

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"go.uber.org/zap"
	"net"
)

type Connection interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Peek(int) ([]byte, error)
	Discard(int) (int, error)

	Close() error
}

type Protocol interface {
	Close() error
	Listen() (Connection, error)
	LocalAddr() net.Addr
	Name() string
}

type Handler interface {
	Receive(Connection) error
}

type Server struct {
	handler  Handler
	protocol Protocol
	logger   *zap.Logger
}

func NewServer(
	rawURL string,
	handler Handler,
	logger *zap.Logger,
) (*Server, error) {
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
		logger:   logger,
	}, nil
}

func (s *Server) ListenAndServe() error {
	s.logger.Info(
		"starting server",
		zap.String("protocol", s.protocol.Name()),
		zap.String("address", s.protocol.LocalAddr().String()),
	)

	protocol := s.protocol
	defer func() {
		s.logger.Info(
			"stopping server",
			zap.String("protocol", s.protocol.Name()),
			zap.String("address", s.protocol.LocalAddr().String()),
		)
		err := protocol.Close()
		if err != nil {
			s.logger.Error(
				"stopping server errored",
				zap.String("protocol", s.protocol.Name()),
				zap.String("address", s.protocol.LocalAddr().String()),
				zap.Error(err),
			)
		}
	}()

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
	go func() {
		s.logger.Info(
			"opening connection",
			zap.String("protocol", s.protocol.Name()),
			zap.String("to", s.protocol.LocalAddr().String()),
			zap.String("from", connection.RemoteAddr().String()),
		)
		err := s.handler.Receive(connection)
		if err != nil {
			s.logger.Error(
				"connection errored",
				zap.String("protocol", s.protocol.Name()),
				zap.String("to", s.protocol.LocalAddr().String()),
				zap.String("from", connection.RemoteAddr().String()),
				zap.Error(err),
			)
		}

		s.logger.Info(
			"closing connection",
			zap.String("protocol", s.protocol.Name()),
			zap.String("to", s.protocol.LocalAddr().String()),
			zap.String("from", connection.RemoteAddr().String()),
		)
	}()

	return nil
}

func (s *Server) Close() error {
	return s.protocol.Close()
}

func (s *Server) LocalAddr() net.Addr {
	return s.protocol.LocalAddr()
}
