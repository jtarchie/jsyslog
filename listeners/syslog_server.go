package listeners

import (
	"github.com/jtarchie/jsyslog/log"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"unsafe"
)

type syslogServer struct {
	*gnet.EventServer
	process ProcessMessage
}

func (u *syslogServer) OnInitComplete(srv gnet.Server) gnet.Action {
	log.Logger.Info(
		"starting server",
		zap.String("protocol", srv.Addr.Network()),
		zap.String("address", srv.Addr.String()),
	)

	return gnet.None
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (u *syslogServer) React(frame []byte, _ gnet.Conn) (out []byte, action gnet.Action) {
	err := u.process(b2s(frame))
	if err != nil {
		return nil, gnet.Close
	}

	return nil, gnet.None
}
