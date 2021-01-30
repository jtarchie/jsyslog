package listeners_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
	"time"
)

var _ = Describe("UDP server", func() {
	Measure("handles many messages", func(b Benchmarker) {
		port, err := listeners.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())

		server, err := listeners.NewUDP("0.0.0.0", port)
		Expect(err).NotTo(HaveOccurred())

		client, err := clients.NewUDP("0.0.0.0", port)
		Expect(err).NotTo(HaveOccurred())

		b.Time("sending messages", func() {
			timer := time.NewTimer(1 * time.Second)
			stopClientServer := make(chan struct{})

			var receivedCounter int32

			go func() {
				err = server.ListenAndServe(func(_ string) error {
					select {
					case <-stopClientServer:
						return fmt.Errorf("time finished")
					default:
						atomic.AddInt32(&receivedCounter, 1)
						return nil
					}
				})
			}()
			var sentCounter int32
			go func() {
				for {
					select {
					case <-stopClientServer:
						return
					default:
						_ = client.WriteString(listeners.PlaceholderValid5424)
						atomic.AddInt32(&sentCounter, 1)
					}
				}
			}()

			<-timer.C
			close(stopClientServer)

			b.RecordValue("number of messages sent", float64(atomic.LoadInt32(&sentCounter)))
			b.RecordValue("number of messages processes", float64(atomic.LoadInt32(&receivedCounter)))
		})

	}, 5)
})
