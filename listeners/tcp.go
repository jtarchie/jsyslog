package listeners

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"github.com/panjf2000/gnet"
	"net"
	"strconv"
	"time"
)

type TCPServer struct {
	address net.TCPAddr
}

type tcpCodec struct {
	gnet.ICodec
}

func (cc *tcpCodec) Decode(c gnet.Conn) ([]byte, error) {
	readLength := 0

	for {
		length, n := c.ReadN(1)
		if length != 1 {
			return nil, fmt.Errorf(
				"could not read message length from TCP server (%s)",
				c.LocalAddr().String(),
			)
		}

		if n[0] == '<' {
			break
		}

		readLength = readLength*10 + int(n[0]-'0')
		_ = c.ShiftN(1)
	}

	actualLength, p := c.ReadN(readLength)

	if actualLength < readLength {
		return nil, fmt.Errorf(
			"could not read from TCP server (%s)",
			c.LocalAddr().String(),
		)
	}

	_ = c.ShiftN(readLength)

	return p, nil
}

func NewTCP(uri *url.URL) (*TCPServer, error) {
	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return nil, fmt.Errorf("could not parse port for TCP server (%s): %w", uri.String(), err)
	}

	address := net.TCPAddr{
		Port: port,
		IP:   net.ParseIP(uri.Hostname()),
	}

	return &TCPServer{
		address: address,
	}, nil
}

func (t *TCPServer) ListenAndServe(process ProcessMessage) error {
	server := &syslogServer{
		process: process,
	}

	err := gnet.Serve(
		server,
		fmt.Sprintf("tcp://%s", t.address.String()),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithCodec(&tcpCodec{}),
		gnet.WithTCPKeepAlive(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf(
			"could not start TCP sever (%s): %w",
			t.address.String(),
			err,
		)
	}

	return nil
}
