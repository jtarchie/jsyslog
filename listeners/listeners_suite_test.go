package listeners_test

import (
	"github.com/jtarchie/jsyslog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestListeners(t *testing.T) {
	log.Logger, _ = zap.NewDevelopment(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewNopCore()
		}),
	)
	log.Logger.Core()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Listeners Suite")
}
