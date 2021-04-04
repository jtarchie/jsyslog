package listeners

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"go.uber.org/zap"
)

type ProcessMessage func(message string) error

type Listener interface {
	ListenAndServe(process ProcessMessage) error
}

func New(rawURL string, logger *zap.Logger) (Listener, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse server from url (%s): %w", rawURL, err)
	}

	switch uri.Scheme {
	case "udp":
		return NewUDP(rawURL, logger)
	case "tcp":
		return NewTCP(rawURL, logger)
	default:
		return nil, fmt.Errorf("could not create server from scheme %q", uri.Scheme)
	}
}
