package listeners

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"github.com/panjf2000/gnet"
	"net"
	"strconv"
)

type UDPServer struct {
	address net.UDPAddr
}

func NewUDP(uri *url.URL) (*UDPServer, error) {
	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return nil, fmt.Errorf("could not parse port for UDP server (%s): %w", uri.String(), err)
	}

	address := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(uri.Hostname()),
	}

	return &UDPServer{
		address: address,
	}, nil
}

func (u *UDPServer) ListenAndServe(process ProcessMessage) error {
	server := &syslogServer{
		process: process,
	}

	err := gnet.Serve(
		server,
		fmt.Sprintf("udp://%s", u.address.String()),
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
	)
	if err != nil {
		return fmt.Errorf(
			"could not start UDP sever (%s): %w",
			u.address.String(),
			err,
		)
	}

	return nil
}
