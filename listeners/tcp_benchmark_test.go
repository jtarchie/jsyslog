package listeners_test

import (
	"fmt"
	"github.com/fgrosse/zaptest"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sync/atomic"
	"time"
)

var _ = Describe("TCP server", func() {
	Measure("handles many messages", func(b Benchmarker) {
		stopClientServer := make(chan struct{})

		port, err := listeners.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())

		logger := zaptest.LoggerWriter(GinkgoWriter)
		var receivedCounter int32
		server, err := listeners.New(
			fmt.Sprintf("tcp://0.0.0.0:%d", port),
			func(msg []byte) error {
				select {
				case <-stopClientServer:
					return fmt.Errorf("time finished")
				default:
					atomic.AddInt32(&receivedCounter, 1)
					return nil
				}
			},
			logger,
		)
		Expect(err).NotTo(HaveOccurred())

		go func() {
			_ = server.ListenAndServe()
		}()

		_, err = clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))
		Expect(err).NotTo(HaveOccurred())

		b.Time("sending messages", func() {
			timer := time.NewTimer(2 * time.Second)

			var sentCounter int32
			for i := 0; i < 5; i++ {
				go func() {
					client, _ := clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))

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
			}

			<-timer.C
			close(stopClientServer)

			b.RecordValue("number of messages sent", float64(atomic.LoadInt32(&sentCounter)))
			b.RecordValue("number of messages processed", float64(atomic.LoadInt32(&receivedCounter)))
		})

	}, 3)
})
