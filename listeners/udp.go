package listeners

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
)

type UDPServer struct {
	logger  *zap.Logger
	process ProcessMessage
	rawURL  string
}

func NewUDP(rawURL string, process ProcessMessage, logger *zap.Logger) (*UDPServer, error) {
	return &UDPServer{
		rawURL:  rawURL,
		logger:  logger,
		process: process,
	}, nil
}

func (u *UDPServer) ListenAndServe() error {
	server := &syslogServer{
		protocol: "udp",
		logger:   u.logger,
		process:  u.process,
	}

	err := gnet.Serve(
		server,
		u.rawURL,
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
	)
	if err != nil {
		return fmt.Errorf("udp server had an issues: %w", err)
	}

	return nil
}
