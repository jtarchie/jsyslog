package listeners_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	"github.com/jtarchie/jsyslog/servers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"log"
	"sync/atomic"
	"time"
)

var _ = Describe("TCP server", func() {
	Measure("handles many messages", func(b Benchmarker) {
		stopClientServer := make(chan struct{})

		port, err := servers.NextReusablePort()
		Expect(err).NotTo(HaveOccurred())

		server, err := listeners.New(fmt.Sprintf("tcp://0.0.0.0:%d", port), zap.NewNop())
		Expect(err).NotTo(HaveOccurred())

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
			log.Printf("server finished: %s", err)
		}()

		_, err = clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))
		Expect(err).NotTo(HaveOccurred())

		b.Time("sending messages", func() {
			timer := time.NewTimer(1 * time.Second)

			var sentCounter int32
			for i := 0; i < 5; i++ {
				go func() {
					client, _ := clients.New(fmt.Sprintf("tcp://0.0.0.0:%d", port))

					for {
						select {
						case <-stopClientServer:
							return
						default:
							_ = client.WriteString(listeners.PlaceholderValid6587)
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

	}, 5)
})
