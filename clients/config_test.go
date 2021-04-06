package clients_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/servers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/jtarchie/jsyslog/clients"
)

var _ = Describe("Configuring Clients", func() {
	for _, protocol := range []string{"tcp", "udp"} {
		protocol := protocol

		When(fmt.Sprintf("the %q client is configured", protocol), func() {
			var (
				server *servers.Server
				port   int
			)

			BeforeEach(func() {
				var err error
				port, err = servers.NextReusablePort()
				Expect(err).NotTo(HaveOccurred())

				server, err = servers.NewServer(
					fmt.Sprintf("%s://0.0.0.0:%d", protocol, port),
					nil,
					zap.NewNop(),
				)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := server.Close()
				Expect(err).NotTo(HaveOccurred())
			})

			It("errors with an unused value", func() {
				_, err := clients.New(fmt.Sprintf("%s://0.0.0.0:%d?unusable-value=1", protocol, port))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nusable-value"))
			})

			It("allows read-deadline as duration", func() {
				_, err := clients.New(fmt.Sprintf("%s://0.0.0.0:%d?read-deadline=1s", protocol, port))
				Expect(err).NotTo(HaveOccurred())

				_, err = clients.New(fmt.Sprintf("%s://0.0.0.0:%d?read-deadline=1", protocol, port))
				Expect(err).To(HaveOccurred())
			})

			It("allows write-deadline as duration", func() {
				_, err := clients.New(fmt.Sprintf("%s://0.0.0.0:%d?write-deadline=1s", protocol, port))
				Expect(err).NotTo(HaveOccurred())

				_, err = clients.New(fmt.Sprintf("%s://0.0.0.0:%d?write-deadline=1", protocol, port))
				Expect(err).To(HaveOccurred())
			})
		})
	}
})
