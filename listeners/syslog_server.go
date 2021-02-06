package listeners

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
)

type syslogServer struct {
	*gnet.EventServer
	process ProcessMessage
}

func (u *syslogServer) OnInitComplete(srv gnet.Server) gnet.Action {
	log.Printf(
		"starting %s",
		fmt.Sprintf("%s://%s",
			srv.Addr.Network(),
			srv.Addr.String(),
		),
	)

	return gnet.None
}

func (u *syslogServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	err := u.process(string(frame))
	if err != nil {
		return nil, gnet.Close
	}

	return nil, gnet.None
}
