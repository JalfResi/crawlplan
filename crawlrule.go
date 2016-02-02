package crawlrate

import (
    "net"
	"time"
    "fmt"
)

type CrawlRule struct {
	Time    time.Duration
	Proxy   net.IP
	Conn    int
	Keyword string
}

func (cr CrawlRule) String() string {
	return fmt.Sprintf("[%ds]\t%s\t%d\t%s\n", int(cr.Time.Seconds()), cr.Proxy.String(), cr.Conn, cr.Keyword)
}