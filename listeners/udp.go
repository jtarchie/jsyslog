package listeners

import (
	"github.com/jtarchie/jsyslog/servers"
	"github.com/jtarchie/jsyslog/url"
)

type UDPServer struct {
	server  *servers.Server
	handler *syslogHandler
}

func NewUDP(uri *url.URL) (*UDPServer, error) {
	handler := &syslogHandler{}
	server, err := servers.NewServer(uri.String(), handler)
	if err != nil {
		return nil, err
	}

	return &UDPServer{
		server:  server,
		handler: handler,
	}, nil
}

func (u *UDPServer) ListenAndServe(process ProcessMessage) error {
	u.handler.process = process

	return u.server.ListenAndServe()
}
