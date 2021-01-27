package payload

import "time"

type Message struct {
	version   int
	priority  int
	hostname  string
	appname   string
	procID    string
	msgID     string
	message   string
	timestamp time.Time
	data      Data
}

func (m *Message) Version() int {
	return m.version
}

func (m *Message) SetVersion(version int) {
	m.version = version
}

func (m *Message) Severity() int {
	return m.priority & 7
}

func (m *Message) Facility() int {
	return m.priority >> 3
}

func (m *Message) Priority() int {
	return m.priority
}

func (m *Message) SetPriority(priority int) {
	m.priority = priority
}

func (m *Message) Hostname() string {
	return m.hostname
}

func (m *Message) SetHostname(hostname string) {
	m.hostname = hostname
}

func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

func (m *Message) SetTimestamp(timestamp time.Time) {
	m.timestamp = timestamp
}

func (m *Message) Appname() string {
	return m.appname
}

func (m *Message) SetAppname(appname string) {
	m.appname = appname
}

func (m *Message) ProcID() string {
	return m.procID
}

func (m *Message) SetProcID(procID string) {
	m.procID = procID
}

func (m *Message) MsgID() string {
	return m.msgID
}

func (m *Message) SetMsgID(msgID string) {
	m.msgID = msgID
}

func (m *Message) Message() string {
	return m.message
}

func (m *Message) SetMessage(message string) {
	m.message = message
}

func (m *Message) StructureData() Data {
	return m.data
}

func (m *Message) AddElement(e *Element) {
	m.data = append(m.data, *e)
}
