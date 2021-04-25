package listeners

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"go.uber.org/zap"
)

type ProcessMessage func(message []byte) error

type Listener interface {
	ListenAndServe() error
}

func New(rawURL string, process ProcessMessage, logger *zap.Logger) (Listener, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse server from url (%s): %w", rawURL, err)
	}

	switch uri.Scheme {
	case "udp":
		return NewUDP(rawURL, process, logger)
	case "tcp":
		return NewTCP(rawURL, process, logger)
	default:
		return nil, fmt.Errorf("could not create server from scheme %q", uri.Scheme)
	}
}
