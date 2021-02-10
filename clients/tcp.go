package clients

import (
	"bufio"
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"math"
	"math/rand"
	"net"
	"time"
)

type TCP struct {
	connection net.Conn
}

func (t *TCP) ReadString() (string, error) {
	reader := bufio.NewReader(t.connection)
	line, _, err := reader.ReadLine()

	return string(line) + "\n", err
}

var _ Client = &TCP{}

func NewTCP(uri *url.URL) (*TCP, error) {
	var connection net.Conn
	err := retry(func() error {
		var err error

		connection, err = net.Dial("tcp", fmt.Sprintf("%s:%s", uri.Hostname(), uri.Port()))
		if err != nil {
			return fmt.Errorf("could not create TCP client (%s:%s): %w", uri.Hostname(), uri.Port(), err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("not successful after retries: %w", err)
	}

	return &TCP{
		connection: connection,
	}, nil
}

func (t *TCP) WriteString(message string) error {
	length, err := t.connection.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("could not write to TCP client (%s): %w", t.connection.LocalAddr(), err)
	}

	if length < len(message) {
		return fmt.Errorf("could not full message to TCP client (%s)", t.connection.LocalAddr())
	}

	return nil
}

func (t *TCP) Close() error {
	err := t.connection.Close()
	if err != nil {
		return fmt.Errorf("could not close TCP client (%s): %w", t.connection.LocalAddr(), err)
	}

	return nil
}

func retry(fn func() error) error {
	retries := 3

	for i := 1; true; i++ {
		err := fn()
		if err != nil {
			if i < retries {
				duration := time.Duration(math.Exp2(float64(i))) * time.Millisecond

				//nolint:gosec
				jitter := time.Duration(rand.Intn(10)) * time.Millisecond

				time.Sleep(duration + jitter)
				continue
			}

			return err
		}
		break
	}

	return nil
}
