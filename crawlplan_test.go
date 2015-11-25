package crawlrate

import (
	"log"
	"net"
	"testing"
	"time"
)

var bottomTests = []struct {
	keywordCount, proxyCount int
	pulse                    *Pulse
	out                      []CrawlRule
}{
	{
		1, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
		},
	},
	{
		2, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(120) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 0, "bingo"},
		},
	},
	{
		2, 1, &Pulse{2, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 1, "bingo"},
		},
	},
	{
		3, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(180) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 0, "bingo"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.3"), 0, "bingo"},
		},
	},
	{
		3, 1, &Pulse{2, time.Duration(60) * time.Second, time.Duration(120) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 1, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.3"), 0, "bingo"},
		},
	},
	{
		3, 1, &Pulse{3, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 1, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.3"), 2, "bingo"},
		},
	},
	{
		15, 3, &Pulse{2, time.Duration(60) * time.Second, time.Duration(180) * time.Second},
		[]CrawlRule{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 1, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 1, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.3"), 0, "bingo"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.3"), 1, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.1"), 1, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 0, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 1, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.3"), 0, "bingo"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.3"), 1, "bingo"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.1"), 0, "bingo"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.2"), 0, "bingo"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.3"), 0, "bingo"},
		},
	},
}

func Test_BottomDistribution(t *testing.T) {
	for k, tt := range bottomTests {
		cp := New(tt.keywordCount, tt.proxyCount, tt.pulse)
		OrderedBy(start, proxy, increasingConnections).Sort(cp)

		log.Printf("test %d %v\n", k, cp)

	}
}
