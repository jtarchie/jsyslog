package main_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/servers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"strings"
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

		outputPath = strings.ReplaceAll(
			writeFile(``),
			`\`,
			`/`,
		)

		bindPort, err = servers.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())
	})

	When("listening on a udp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "forwarder",
				"--from", fmt.Sprintf("udp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("file://%s", outputPath),
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting server`))
			writeUDP(bindPort, validUDPMessage)

			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validUDPMessage)),
			)
		})
	})

	When("listening on a tcp://", func() {
		It("listens and writes syslog messages to a file", func() {
			command := exec.Command(binPath, "forwarder",
				"--from", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("file://%s", outputPath),
			)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer session.Kill()
			Eventually(session.Err).Should(gbytes.Say(`starting server`))
			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validTCPMessage), validTCPMessage))
			Eventually(session.Err).Should(gbytes.Say(`opening connection`))

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
			Eventually(session.Err).Should(gbytes.Say(`starting server`))

			writeTCP(bindPort, fmt.Sprintf("%d%s", len(validTCPMessage), validTCPMessage))
			writeUDP(bindPort, validUDPMessage)
			Eventually(session.Err).Should(gbytes.Say(`opening connection`))
			Eventually(session.Err).Should(gbytes.Say(`closing connection`))

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
			Eventually(secondSession.Err).Should(gbytes.Say(`starting server`))

			firstSession, err := gexec.Start(exec.Command(binPath,
				"forwarder",
				"--from", fmt.Sprintf("udp://0.0.0.0:%d", bindPort),
				"--to", fmt.Sprintf("tcp://0.0.0.0:%d", bindPort),
			), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			defer firstSession.Kill()
			Eventually(firstSession.Err).Should(gbytes.Say(`starting server`))

			writeUDP(bindPort, validUDPMessage)

			Eventually(readFile(outputPath)).Should(
				ContainSubstring(fmt.Sprintf("%s\n", validUDPMessage)),
			)
		})
	})
})
