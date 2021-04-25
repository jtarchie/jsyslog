package listeners

import (
	"github.com/panjf2000/gnet"
	"go.uber.org/zap"
)

type syslogServer struct {
	*gnet.EventServer
	process  ProcessMessage
	protocol string
	logger   *zap.Logger
}

func (t *syslogServer) OnInitComplete(server gnet.Server) gnet.Action {
	t.logger.Info(
		"started server",
		zap.String("to", server.Addr.String()),
		zap.String("protocol", t.protocol),
	)

	return gnet.None
}

func (t *syslogServer) OnShutdown(server gnet.Server) {
	t.logger.Info(
		"stopping server",
		zap.String("to", server.Addr.String()),
		zap.String("protocol", t.protocol),
	)
}

func (t *syslogServer) OnOpened(c gnet.Conn) ([]byte, gnet.Action) {
	t.logger.Info(
		"accepting connection",
		zap.String("to", c.LocalAddr().String()),
		zap.String("from", c.RemoteAddr().String()),
		zap.String("protocol", t.protocol),
	)

	return nil, gnet.None
}

func (t *syslogServer) OnClosed(c gnet.Conn, err error) gnet.Action {
	t.logger.Info(
		"closing connection",
		zap.String("to", c.LocalAddr().String()),
		zap.String("from", c.RemoteAddr().String()),
		zap.String("protocol", t.protocol),
		zap.Error(err),
	)

	return gnet.None
}

func (t *syslogServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	err := t.process(frame)
	if err != nil {
		t.logger.Error(
			"could not process message",
			zap.String("to", c.LocalAddr().String()),
			zap.String("from", c.RemoteAddr().String()),
			zap.String("protocol", t.protocol),
			zap.Error(err),
		)
	}

	return nil, gnet.None
}
