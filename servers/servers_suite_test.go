package servers_test

import (
	"bufio"
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/servers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"testing"
)

func TestServers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Servers Suite")
}

type echoHandler struct{}

func (e echoHandler) Receive(connection servers.Connection) error {
	reader := bufio.NewReader(connection)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	_, err = connection.Write(append(line, []byte("\n")...))
	return err
}

var _ servers.Handler = &echoHandler{}

var _ = Describe("When creating servers", func() {
	When("a UDP server echos back", func() {
		It("captures the echo", func() {
			port, err := servers.NextReusablePort()
			Expect(err).NotTo(HaveOccurred())

			server, err := servers.NewServer(
				fmt.Sprintf("udp://0.0.0.0:%d", port),
				&echoHandler{},
				zap.NewNop(),
			)
			Expect(err).NotTo(HaveOccurred())
			defer server.Close()

			go func() {
				_ = server.ListenAndServe()
			}()
			response := writeUDP(port, "Hello, World!\n")
			Expect(response).To(Equal("Hello, World!\n"))
		})
	})

	When("a TCP server echos back", func() {
		It("captures the echo", func() {
			port, err := servers.NextReusablePort()
			Expect(err).NotTo(HaveOccurred())

			server, err := servers.NewServer(
				fmt.Sprintf("tcp://0.0.0.0:%d", port),
				&echoHandler{},
				zap.NewNop(),
			)
			Expect(err).NotTo(HaveOccurred())
			defer server.Close()

			go func() {
				_ = server.ListenAndServe()
			}()
			response := writeTCP(port, "Hello, World!\n")
			Expect(response).To(Equal("14Hello, World!\n"))
		})
	})

	When("TCP and UDP bind to the same port", func() {
		It("allows it to happen", func() {
			port, err := servers.NextReusablePort()
			Expect(err).NotTo(HaveOccurred())

			tcpServer, err := servers.NewServer(
				fmt.Sprintf("tcp://0.0.0.0:%d", port),
				&echoHandler{},
				zap.NewNop(),
			)
			Expect(err).NotTo(HaveOccurred())
			defer tcpServer.Close()

			udpServer, err := servers.NewServer(
				fmt.Sprintf("udp://0.0.0.0:%d", port),
				&echoHandler{},
				zap.NewNop(),
			)
			Expect(err).NotTo(HaveOccurred())
			defer udpServer.Close()

			go func() {
				_ = tcpServer.ListenAndServe()
			}()

			go func() {
				_ = udpServer.ListenAndServe()
			}()

			tcpResponse := writeTCP(port, "Hello, World!\n")
			Expect(tcpResponse).To(Equal("14Hello, World!\n"))

			udpResponse := writeUDP(port, "Hello, World!\n")
			Expect(udpResponse).To(Equal("Hello, World!\n"))
		})
	})
})

func writeUDP(port int, message string) string {
	client, err := clients.New(fmt.Sprintf("udp://0.0.0.0:%d?read-timeout=1s", port))
	Expect(err).NotTo(HaveOccurred())
	Expect(client).To(BeAssignableToTypeOf(&clients.UDP{}))

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	response, err := client.ReadString()
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())

	return response
}

func writeTCP(port int, message string) string {
	client, err := clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))
	Expect(err).NotTo(HaveOccurred())
	Expect(client).To(BeAssignableToTypeOf(&clients.SyslogTCP{}))

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	response, err := client.ReadString()
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())

	return response
}
