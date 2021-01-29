package listeners_test

import (
	"github.com/jtarchie/jsyslog/clients"
	"github.com/jtarchie/jsyslog/listeners"
	"sync/atomic"
	"testing"
)

func BenchmarkUDPServer(b *testing.B) {
	port, err := listeners.NextReusablePort()
	if err != nil {
		b.Fatalf("could not find a port: %s", err)
	}

	server, err := listeners.NewUDP("0.0.0.0", port)
	if err != nil {
		b.Fatalf("could not start server: %s", err)
	}

	var counter int32
	go server.ListenAndServe(func(_ string) error {
		atomic.AddInt32(&counter, 1)
		return nil
	})

	client, err := clients.NewUDP("0.0.0.0", port)
	if err != nil {
		b.Fatalf("could not create client: %s", err)
	}

	b.Run("sending messages", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = client.WriteString(listeners.PlaceholderValid5424)
		}
	})
}
