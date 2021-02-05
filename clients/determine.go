package clients

import (
	"fmt"
	"net/url"
)

type Client interface {
	WriteString(string) error
	Close() error
}

func New(rawURL string) (Client, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse client from url (%s): %w", rawURL, err)
	}

	switch uri.Scheme {
	case "udp":
		return NewUDP(uri)
	case "tcp":
		return NewTCP(uri)
	default:
		return nil, fmt.Errorf("could not create client from scheme %q", uri.Scheme)
	}
}
