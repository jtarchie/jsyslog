package main_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"io/ioutil"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForwarder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSyslog Suite")
}

const validUDPMessage = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' UDP`
const validTCPMessage = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' TCP`

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
	Expect(client).To(BeAssignableToTypeOf(&clients.UDP{}))

	err = client.WriteString(message)
	Expect(err).NotTo(HaveOccurred())

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())
}

func writeTCP(port int, message string) {
	client, err := clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))
	Expect(err).NotTo(HaveOccurred())
	Expect(client).To(BeAssignableToTypeOf(&clients.TCP{}))

	for i := 0; i < 10; i++ {
		err = client.WriteString(message)
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("failed on message %d", i))
		time.Sleep(10 * time.Millisecond)
	}

	err = client.Close()
	Expect(err).NotTo(HaveOccurred())
}
