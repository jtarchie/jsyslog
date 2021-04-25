package listeners

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
)

type TCPServer struct {
	logger  *zap.Logger
	process ProcessMessage
	rawURL  string
}

type tcpSyslogCodec struct {
	gnet.ICodec
	logger *zap.Logger
}

func (s *tcpSyslogCodec) Decode(c gnet.Conn) ([]byte, error) {
	readLength := 0
	readDigits := 0

	for {

		length, n := c.ReadN(readDigits + 1)
		if length != readDigits+1 {
			return nil, fmt.Errorf(
				"could not read message length from TCP server (%s)",
				c.LocalAddr().String(),
			)
		}

		if n[readDigits] == '<' {
			break
		}

		readLength = readLength*10 + int(n[readDigits]-'0')
		readDigits++
	}

	expectReadLength := readLength + readDigits
	actualLength, p := c.ReadN(expectReadLength)

	if actualLength < expectReadLength {
		return nil, fmt.Errorf(
			"could not read from TCP server (%s)",
			c.LocalAddr().String(),
		)
	}

	_ = c.ShiftN(expectReadLength)

	message := p[readDigits:]
	return message, nil
}

func NewTCP(rawURL string, process ProcessMessage, logger *zap.Logger) (*TCPServer, error) {
	return &TCPServer{
		rawURL:  rawURL,
		logger:  logger,
		process: process,
	}, nil
}

func (t *TCPServer) ListenAndServe() error {
	server := &syslogServer{
		protocol: "tcp",
		logger:   t.logger,
		process:  t.process,
	}

	err := gnet.Serve(
		server,
		t.rawURL,
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithCodec(&tcpSyslogCodec{
			logger: t.logger,
		}),
	)
	if err != nil {
		return fmt.Errorf("tcp server had an issues: %w", err)
	}

	return nil
}
