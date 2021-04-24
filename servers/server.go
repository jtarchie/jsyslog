package servers

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
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
	handler          Handler
	logger           *zap.Logger
	protocol         Protocol
	totalConnections uint64
	worker           *Worker
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
		worker:   NewWorker(15, logger),
	}, nil
}

func (s *Server) ListenAndServe() error {
	logger := s.logger
	protocol := s.protocol
	worker := s.worker

	worker.Start()
	logger.Info(
		"starting server",
		zap.String("protocol", protocol.Name()),
		zap.String("address", protocol.LocalAddr().String()),
	)

	for {
		connection, err := protocol.Listen()
		if err != nil {
			return fmt.Errorf("could not listen for connection (%s): %w", protocol.LocalAddr().String(), err)
		}

		id := atomic.AddUint64(&s.totalConnections, 1)
		worker.Run(func(workerID int) error {
			logger.Info(
				"opening connection",
				zap.String("protocol", protocol.Name()),
				zap.String("address", protocol.LocalAddr().String()),
				zap.Uint64("connectionID", id),
				zap.Int("workerID", workerID),
			)

			defer logger.Info(
				"closing connection",
				zap.String("protocol", protocol.Name()),
				zap.String("address", protocol.LocalAddr().String()),
				zap.Uint64("connectionID", id),
				zap.Int("workerID", workerID),
			)
			return s.handler.Receive(connection)
		})
	}
}

func (s *Server) Close() error {
	logger := s.logger
	protocol := s.protocol
	worker := s.worker

	defer worker.Stop()

	logger.Info(
		"stopping server",
		zap.String("protocol", protocol.Name()),
		zap.String("address", protocol.LocalAddr().String()),
	)
	err := protocol.Close()
	if err != nil {
		logger.Error(
			"stopping server errored",
			zap.String("protocol", protocol.Name()),
			zap.String("address", protocol.LocalAddr().String()),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *Server) LocalAddr() net.Addr {
	return s.protocol.LocalAddr()
}
