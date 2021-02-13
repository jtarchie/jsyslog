package listeners

import (
	"github.com/jtarchie/jsyslog/log"
	"github.com/jtarchie/jsyslog/servers"
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
	"io/ioutil"
	"unsafe"
)

type syslogServer struct {
	*gnet.EventServer
	process ProcessMessage
}

func (u *syslogServer) OnOpened(c gnet.Conn) ([]byte, gnet.Action) {
	log.Logger.Info(
		"opening connection",
		zap.String("local", c.LocalAddr().String()),
		zap.String("protocol", c.LocalAddr().Network()),
		zap.String("remote", c.RemoteAddr().String()),
	)

	return nil, gnet.None
}

func (u *syslogServer) OnClosed(c gnet.Conn, _ error) gnet.Action {
	log.Logger.Info(
		"closing connection",
		zap.String("local", c.LocalAddr().String()),
		zap.String("protocol", c.LocalAddr().Network()),
		zap.String("remote", c.RemoteAddr().String()),
	)

	return gnet.None
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

func (u *syslogServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	log.Logger.Debug(
		"received message",
		zap.String("local", c.LocalAddr().String()),
		zap.String("protocol", c.LocalAddr().Network()),
		zap.String("remote", c.RemoteAddr().String()),
	)

	err := u.process(b2s(frame))
	if err != nil {
		return nil, gnet.Close
	}

	return nil, gnet.None
}

type syslogHandler struct {
	process ProcessMessage
}

var _ servers.Handler = &syslogHandler{}

func (s syslogHandler) Receive(connection servers.Connection) error {
	message, err := ioutil.ReadAll(connection)
	if err != nil {
		return err
	}

	err = s.process(string(message))
	if err != nil {
		return err
	}

	return nil
}
