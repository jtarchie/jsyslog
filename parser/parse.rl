package parser

import (
    "github.com/jtarchie/jsyslog/payload"
    "fmt"
    "time"
)

%%{
    machine syslog;
    write data;
}%%

var UnknownError = fmt.Errorf("unknown parsing error occurred, check offset for when it failed")

func Parse(data string) (*payload.Message, int, error) {
    // required for ragel state management
    cs, p, pe, eof := 0, 0, len(data), len(data)

    // keep track of where a state entered
    mark := 0
    // allow a partial parsing
    message := &payload.Message{}
    var (
        // keep track of timezone location when parsing
        location *time.Location
        // keep track of nanoseconds from the timestamp
        nanosecond int
        // create a buffer for capturing sd-values
        buffer []byte
        element *payload.Element
        paramName string
    )

%%{
    action mark { mark = p }

    action tcp_len { pe, eof = atoi(data[mark:p]) + (p-mark) + 1, atoi(data[mark:p]) + (p-mark) + 1 }

    action version  { message.SetVersion(atoi(data[mark:p])) }
    action priority { message.SetPriority(atoi(data[mark:p])) }
    action hostname { message.SetHostname(data[mark:p]) }
    action appname  { message.SetAppname(data[mark:p]) }
    action procid   { message.SetProcID(data[mark:p]) }
    action msgid    { message.SetMsgID(data[mark:p]) }
    action message  { message.SetMessage(data[mark:p]) }

    action escaped    {
      buffer = append(buffer, data[mark:p-2]...)
      buffer = append(buffer, data[p-1])
      mark = p
    }
    action paramname  { paramName = data[mark:p] }
    action paramvalue {
      buffer = append(buffer, data[mark:p]...)
      element.AddProperty(payload.NewProperty(paramName, string(buffer)))
      buffer = buffer[:0]
    }
    action sdid      {
      if element != nil {
        message.AddElement(element)
      }
      element = payload.NewElement(data[mark:p])
    }

    action timestamp {
      location = time.UTC
      if data[mark+19] == '.' {
        offset := 1

        if data[p-1] != 'Z' {
          offset = 6
          dir := 1
          if data[p-6] == '-' {
            dir = -1
          }

          location = time.FixedZone(
            "",
            dir * (atoi2(data[p-5:p-3]) * 3600 + atoi(data[p-2:p]) * 60),
          )
        }
        nbytes := ( p - offset - 1 ) - ( mark + 19 )
        i := mark + 20
        first:
          if i < p-offset {
            nanosecond = nanosecond*10 + int(data[i]-'0')
            i++
            goto first
          }
        i = 0
        second:
          if i < 9-nbytes {
            nanosecond *= 10
            i++
            goto second
          }
      }

      message.SetTimestamp(time.Date(
        atoi4(data[mark:mark+4]),
        time.Month(atoi2(data[mark+5:mark+7])),
        atoi2(data[mark+8:mark+10]),
        atoi2(data[mark+11:mark+13]),
        atoi2(data[mark+14:mark+16]),
        atoi2(data[mark+17:mark+19]),
        nanosecond,
        location,
      ).UTC())
    }

    include syslog_rfc5424 "syslog_rfc5424.rl";
    main := tcp_syslog_msg | syslog_msg;

    write init;
    write exec;
}%%
    if element != nil {
      message.AddElement(element)
    }

    if cs < syslog_first_final {
        return message, p, UnknownError
    }

    return message, p, nil
}
