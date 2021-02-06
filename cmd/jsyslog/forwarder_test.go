package main_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/listeners"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("When forwarding messages", func() {
	var (
		binPath    string
		outputPath string
		bindPort   int
	)

	BeforeEach(func() {
		var err error
		binPath, err = gexec.Build("github.com/jtarchie/jsyslog/cmd/jsyslog")
		Expect(err).NotTo(HaveOccurred())

		outputPath = writeFile(``)

		bindPort, err = listeners.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())
	})

	When("listening on a udp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "forwarder", "--from", fmt.Sprintf("udp://0.0.0.0:%d", bindPort), "--to", fmt.Sprintf("file://%s", outputPath))
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
			command := exec.Command(binPath, "forwarder", "--from", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort), "--to", fmt.Sprintf("file://%s", outputPath))
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting tcp://0.0.0.0`))

			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validMessage), validMessage))
			Eventually(readFile(outputPath)).Should(Equal(validMessage))
		})
	})
})
