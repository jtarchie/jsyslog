package main_test

import (
	"fmt"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"net"
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
	It("listens and writes syslog messages to a file", func() {
		path, err := gexec.Build("github.com/jtarchie/jsyslog/cmd/storer")
		Expect(err).NotTo(HaveOccurred())

		filename := writeFile(``)
		port := getFreePort()

		command := exec.Command(path, "--listen", fmt.Sprintf("udp://0.0.0.0:%d", port), "--file", filename)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		defer session.Kill()
		Eventually(session.Err).Should(gbytes.Say(`starting UDP server`))

		writeUDP(port, validMessage)
		Eventually(readFile(filename)).Should(ContainSubstring(validMessage))
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

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	Expect(err).NotTo(HaveOccurred())

	l, err := net.ListenTCP("tcp", addr)
	Expect(err).NotTo(HaveOccurred())

	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func readFile(filename string) func() string {
	return func() string {
		contents, err := ioutil.ReadFile(filename)
		Expect(err).NotTo(HaveOccurred())

		return string(contents)
	}
}

func writeUDP(port int, message string) {
	conn, err := net.Dial("udp", fmt.Sprintf("0.0.0.0:%d", port))
	Expect(err).NotTo(HaveOccurred())

	length, err := fmt.Fprintf(conn, message)
	Expect(err).NotTo(HaveOccurred())
	Expect(length).To(Equal(len(message)))

	err = conn.Close()
	Expect(err).NotTo(HaveOccurred())
}
