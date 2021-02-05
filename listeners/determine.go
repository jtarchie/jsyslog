package listeners

import (
	"fmt"
	"net/url"
)

type ProcessMessage func(message string) error

type Listener interface {
	ListenAndServe(process ProcessMessage) error
}

func New(rawURL string) (Listener, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse server from url (%s): %w", rawURL, err)
	}

	switch uri.Scheme {
	case "udp":
		return NewUDP(uri)
	case "tcp":
		return NewTCP(uri)
	default:
		return nil, fmt.Errorf("could not create server from scheme %q", uri.Scheme)
	}
}
