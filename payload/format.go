package payload

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Format int

const (
	RFC5424 Format = iota
)

func (m *Message) String() string {
	return m.Format(RFC5424)
}

func (m *Message) Format(f Format) string {
	if f == RFC5424 {
		return FormatRFC5424(m)
	}

	return ""
}

var paramValueReplacer = strings.NewReplacer(`"`, `\"`, `]`, `\]`, `\`, `\\`)

func FormatRFC5424(m *Message) string {
	var buffer bytes.Buffer
	buffer.WriteString("<")
	buffer.WriteString(fmt.Sprintf("%d", m.Priority()))
	buffer.WriteString(">")
	buffer.WriteString(fmt.Sprintf("%d", m.Version()))
	buffer.WriteString(" ")
	if m.timestamp.IsZero() {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.timestamp.Format(time.RFC3339Nano))
	}
	buffer.WriteString(" ")
	if m.Hostname() == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.Hostname())
	}
	buffer.WriteString(" ")
	if m.Appname() == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.Appname())
	}
	buffer.WriteString(" ")
	if m.ProcID() == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.ProcID())
	}
	buffer.WriteString(" ")
	if m.MsgID() == "" {
		buffer.WriteString("-")
	} else {
		buffer.WriteString(m.MsgID())
	}
	buffer.WriteString(" ")

	if len(m.StructureData()) == 0 {
		buffer.WriteString("-")
	} else {
		for _, element := range m.StructureData() {
			if element.ID() == "" {
				buffer.WriteString("-")
			} else {
				buffer.WriteString("[")
				buffer.WriteString(element.ID())

				for _, property := range element.Properties() {
					buffer.WriteString(" ")
					buffer.WriteString(property.Key())
					buffer.WriteString(`="`)
					buffer.WriteString(paramValueReplacer.Replace(property.Value()))
					buffer.WriteString(`"`)
				}

				buffer.WriteString("]")
			}
		}
	}

	if m.Message() != "" {
		buffer.WriteString(" ")
		buffer.WriteString(m.Message())
	}
	return buffer.String()
}
