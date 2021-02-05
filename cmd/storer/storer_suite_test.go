package main_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForwarder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storer Suite")
}

const validMessage = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`

var _ = Describe("When it executes", func() {
	var (
		binPath    string
		outputPath string
		bindPort   int
	)

	BeforeEach(func() {
		var err error
		binPath, err = gexec.Build("github.com/jtarchie/jsyslog/cmd/storer")
		Expect(err).NotTo(HaveOccurred())

		outputPath = writeFile(``)

		bindPort, err = listeners.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())
	})

	When("listening on a udp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "--listen", fmt.Sprintf("udp://0.0.0.0:%d", bindPort), "--file", outputPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting udp://0.0.0.0`))

			writeUDP(bindPort, validMessage)
			Eventually(readFile(outputPath)).Should(Equal(validMessage))
		})
	})

	When("listening on a tcp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "--listen", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort), "--file", outputPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting tcp://0.0.0.0`))

			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validMessage), validMessage))
			Eventually(readFile(outputPath)).Should(Equal(validMessage))
		})
	})
})

func writeFile(contents string) string {
	file, err := ioutil.TempFile("", "")
	Expect(err).NotTo(HaveOccurred())

	_, err = file.WriteString(contents)
	Expect(err).NotTo(HaveOccurred())

	err = file.Close()
	Expect(err).NotTo(HaveOccurred())

	return file.Name()
}

func readFile(filename string) func() string {
	return func() string {
		contents, err := ioutil.ReadFile(filename)
		Expect(err).NotTo(HaveOccurred())

		return string(contents)
	}
}

func writeUDP(port int, message string) {
	client, err := clients.New(fmt.Sprintf("udp://0.0.0.0:%d", port))
	Expect(err).NotTo(HaveOccurred())
	Expect(client).To(BeAssignableToTypeOf(&clients.UDPClient{}))

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())
}

func writeTCP(port int, message string) {
	client, err := clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))
	Expect(err).NotTo(HaveOccurred())
	Expect(client).To(BeAssignableToTypeOf(&clients.TCPClient{}))

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())
}
