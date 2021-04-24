package servers_test

import (
	"fmt"
	"github.com/fgrosse/zaptest"
	"github.com/jtarchie/jsyslog/servers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var _ = Describe("Worker", func() {
	var logger *zap.Logger

	BeforeEach(func() {
		logger = zaptest.LoggerWriter(GinkgoWriter)
	})

	When("with only one worker", func() {
		It("can only run one thing at a time", func() {
			counter := int64(0)
			worker := servers.NewWorker(1, logger)
			defer worker.Stop()

			worker.Start()
			worker.Run(func(_ int) error {
				atomic.AddInt64(&counter, 1)
				return nil
			})

			Eventually(func() int64 {
				return atomic.LoadInt64(&counter)
			}).Should(BeNumerically("==", 1))

			Consistently(func() int64 {
				return atomic.LoadInt64(&counter)
			}).Should(BeNumerically("==", 1))
		})

		It("only starts one goroutine", func() {
			worker := servers.NewWorker(1, logger)

			numRoutines := runtime.NumGoroutine()
			worker.Start()
			Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+1))
			Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+1))

			worker.Start()
			Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+1))
			Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+1))

			worker.Stop()
			Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))
			Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))

			worker.Stop()
			Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))
			Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))
		})

		It("allows errors to be processed", func() {
			counter := int64(0)
			worker := servers.NewWorker(
				1,
				logger,
				servers.WithErrorFunc(func(err error) {
					atomic.AddInt64(&counter, 1)
				}),
			)

			defer worker.Stop()

			worker.Start()
			worker.Run(func(_ int) error {
				return fmt.Errorf("some error occurred")
			})

			Eventually(func() int64 {
				return atomic.LoadInt64(&counter)
			}).Should(BeNumerically("==", 1))

			Consistently(func() int64 {
				return atomic.LoadInt64(&counter)
			}).Should(BeNumerically("==", 1))
		})
	})

	When("more than one worker", func() {
		for _, numberOfWorkers := range []int{2, 5, 10} {
			It("can only run one thing at a time per worker", func() {
				Expect(numberOfWorkers).To(BeNumerically(">", 0))

				counter := int64(0)
				worker := servers.NewWorker(uint(numberOfWorkers), logger)
				activeWorkers := &sync.Map{}

				defer worker.Stop()

				worker.Start()

				for i := 0; i < numberOfWorkers; i++ {
					worker.Run(func(id int) error {
						activeWorkers.Store(id, 1)
						atomic.AddInt64(&counter, 1)

						time.Sleep(10 * time.Millisecond)
						return nil
					})
				}

				Eventually(func() int64 {
					return atomic.LoadInt64(&counter)
				}).Should(BeNumerically("==", numberOfWorkers))

				Consistently(func() int64 {
					return atomic.LoadInt64(&counter)
				}).Should(BeNumerically("==", numberOfWorkers))

				totalWorkers := 0
				activeWorkers.Range(func(_, _ interface{}) bool {
					totalWorkers++
					return true
				})
				Expect(totalWorkers).To(Equal(numberOfWorkers))
			})

			It(fmt.Sprintf("only starts %d goroutine", numberOfWorkers), func() {
				worker := servers.NewWorker(uint(numberOfWorkers), logger)

				numRoutines := runtime.NumGoroutine()
				Consistently(runtime.NumGoroutine).Should(BeNumerically("~", numRoutines, 1))

				worker.Start()

				Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+numberOfWorkers))
				Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines+numberOfWorkers))

				worker.Stop()

				Eventually(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))
				Consistently(runtime.NumGoroutine).Should(BeNumerically("==", numRoutines))
			})
		}
	})
})
