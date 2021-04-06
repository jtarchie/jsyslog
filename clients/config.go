package clients

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"time"
)

type configuration struct {
	readDeadline  time.Duration
	writeDeadline time.Duration
}

func contains(haystack []string, needle string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

func newConfig(uri *url.URL) (*configuration, error) {
	validParams := []string{"read-deadline", "write-deadline"}

	for name := range uri.Query() {
		if !contains(validParams, name) {
			return nil, fmt.Errorf("cannot configure %q on the connection", name)
		}
	}

	config := &configuration{
		readDeadline:  24 * time.Hour,
		writeDeadline: 24 * time.Hour,
	}

	if timeout := uri.Query().Get("read-deadline"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse read deadline duration for connection: %w",
				err,
			)
		}
		config.readDeadline = duration
	}

	if timeout := uri.Query().Get("write-deadline"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse write deadline duration for connection: %w",
				err,
			)
		}
		config.writeDeadline = duration
	}

	return config, nil
}
