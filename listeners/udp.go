package listeners

import (
	"github.com/jtarchie/jsyslog/servers"
	"io/ioutil"
)

type UDPServer struct {
	server  *servers.Server
	process ProcessMessage
}

func (u *UDPServer) Receive(connection servers.Connection) error {
	message, err := ioutil.ReadAll(connection)
	if err != nil {
		return err
	}

	err = u.process(string(message))
	if err != nil {
		return err
	}

	return nil
}

func NewUDP(rawURL string) (*UDPServer, error) {
	handler := &UDPServer{}

	server, err := servers.NewServer(rawURL, handler)
	if err != nil {
		return nil, err
	}

	handler.server = server

	return handler, nil
}

func (u *UDPServer) ListenAndServe(process ProcessMessage) error {
	u.process = process

	return u.server.ListenAndServe()
}
