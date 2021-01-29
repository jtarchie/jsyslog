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
	It("listens and writes syslog messages to a file", func() {
		path, err := gexec.Build("github.com/jtarchie/jsyslog/cmd/storer")
		Expect(err).NotTo(HaveOccurred())

		filename := writeFile(``)

		port, err := listeners.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(path, "--listen", fmt.Sprintf("udp://0.0.0.0:%d", port), "--file", filename)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		defer session.Kill()
		Eventually(session.Err).Should(gbytes.Say(`starting udp://0.0.0.0`))

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

func readFile(filename string) func() string {
	return func() string {
		contents, err := ioutil.ReadFile(filename)
		Expect(err).NotTo(HaveOccurred())

		return string(contents)
	}
}

func writeUDP(port int, message string) {
	client, err := clients.NewUDP("0.0.0.0", port)
	Expect(err).NotTo(HaveOccurred())

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())
}
