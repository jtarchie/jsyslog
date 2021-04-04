package clients

import (
	"fmt"
	"github.com/jtarchie/jsyslog/url"
	"strconv"
)

type SyslogTCP struct {
	*TCP
}

func NewSyslogTCP(uri *url.URL) (*SyslogTCP, error) {
	client, err := NewTCP(uri)
	if err != nil {
		return nil, fmt.Errorf("could not create syslog TCP client: %w", err)
	}

	return &SyslogTCP{
		TCP: client,
	}, nil
}

func (t *SyslogTCP) WriteString(message string) error {
	_, err := t.connection.Write([]byte(strconv.Itoa(len(message))))
	if err != nil {
		return fmt.Errorf("could not write message length to TCP client (%s): %w", t.connection.LocalAddr(), err)
	}

	return t.TCP.WriteString(message)
}
