package listeners

import (
	"fmt"
	"github.com/jtarchie/jsyslog/servers"
	"github.com/jtarchie/jsyslog/url"
)

type TCPServer struct {
	server  *servers.Server
	process ProcessMessage
}

func (t *TCPServer) Receive(connection servers.Connection) error {
	readLength := 0

	for {
		n, err := connection.Peek(1)

		if err != nil {
			return fmt.Errorf(
				"could not read message from TCP server (%s): %w",
				t.server.LocalAddr(),
				err,
			)
		}

		if len(n) != 1 {
			return fmt.Errorf(
				"could not read message length from TCP server (%s)",
				t.server.LocalAddr(),
			)
		}

		if n[0] == '<' && readLength > 0 {
			goto readMessage
		}

		if '0' <= n[0] && n[0] <= '9' {
			readLength = readLength*10 + int(n[0]-'0')
		} else {
			readLength = 0
		}

		_, _ = connection.Discard(1)
		continue

	readMessage:
		p, _ := connection.Peek(readLength)

		actualLength := len(p)

		if actualLength < readLength {
			return fmt.Errorf(
				"could not read from TCP server (%s)",
				t.server.LocalAddr(),
			)
		}

		_, _ = connection.Discard(actualLength)

		err = t.process(string(p))
		if err != nil {
			return err
		}

		readLength = 0
	}
}

func NewTCP(uri *url.URL) (*TCPServer, error) {
	handler := &TCPServer{}

	server, err := servers.NewServer(uri.String(), handler)
	if err != nil {
		return nil, err
	}

	handler.server = server

	return handler, nil
}

func (t *TCPServer) ListenAndServe(process ProcessMessage) error {
	t.process = process

	return t.server.ListenAndServe()
}
