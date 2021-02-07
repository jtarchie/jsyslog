package main_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/listeners"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"time"
)

var _ = Describe("When forwarding messages", func() {
	var (
		binPath    string
		outputPath string
		bindPort   int
	)

	BeforeEach(func() {
		var err error
		binPath, err = gexec.Build("github.com/jtarchie/jsyslog")
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

			writeUDP(bindPort, validUDPMessage)
			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validUDPMessage)),
			)
		})
	})

	When("listening on a tcp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "forwarder", "--from", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort), "--to", fmt.Sprintf("file://%s", outputPath))
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting tcp://0.0.0.0`))

			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validTCPMessage), validTCPMessage))
			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validTCPMessage)),
			)
		})
	})

	When("listening for udp and tcp on the same port", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "forwarder",
				"--from", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort),
				"--from", fmt.Sprintf("udp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("file://%s", outputPath),
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting udp://0.0.0.0`))

			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validTCPMessage), validTCPMessage))
			writeUDP(bindPort, validUDPMessage)

			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validUDPMessage)),
			)
			Eventually(readFile(outputPath), 2*time.Second).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validTCPMessage)),
			)
		})
	})

	When("forwarding UDP to TCP to file", func() {
		It("listens and writes syslog messages to a file", func() {
			secondSession, err := gexec.Start(exec.Command(binPath,
				"forwarder",
				"--from", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("file://%s", outputPath),
			), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer secondSession.Kill()
			Eventually(secondSession.Err).Should(gbytes.Say(`starting tcp://0.0.0.0`))

			firstSession, err := gexec.Start(exec.Command(binPath,
				"forwarder",
				"--from", fmt.Sprintf("udp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort),
			), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer firstSession.Kill()
			Eventually(firstSession.Err).Should(gbytes.Say(`starting udp://0.0.0.0`))

			writeUDP(bindPort, validUDPMessage)

			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validUDPMessage)),
			)
		})
	})
})
