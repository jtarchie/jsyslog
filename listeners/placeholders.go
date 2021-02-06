package listeners

import "fmt"

const (
	PlaceholderValid5424 = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`
)

var PlaceholderValid6587 = fmt.Sprintf("%d%s", len(PlaceholderValid5424), PlaceholderValid5424)
