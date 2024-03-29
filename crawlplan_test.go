package crawlrate

import (
	"net"
	"testing"
	"time"
	"fmt"
)

var bottomTests = []struct {
	keywordCount, proxyCount int
	pulse                    *Pulse
	out                      CrawlPlan
}{
	{
		1, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
		},
	},
	{
		2, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(120) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-1"},
		},
	},
	{
		2, 1, &Pulse{2, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-1"},
		},
	},
	{
		3, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(180) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-1"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-2"},
		},
	},
	{
		3, 1, &Pulse{2, time.Duration(60) * time.Second, time.Duration(120) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-1"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-2"},
		},
	},
	{
		3, 1, &Pulse{3, time.Duration(60) * time.Second, time.Duration(60) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-1"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 2, "keyword-2"},
		},
	},
	{
		15, 3, &Pulse{2, time.Duration(60) * time.Second, time.Duration(180) * time.Second},
		CrawlPlan{
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-1"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 0, "keyword-2"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.1"), 1, "keyword-3"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 0, "keyword-4"},
			{time.Duration(0) * time.Second, net.ParseIP("127.0.0.2"), 1, "keyword-5"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-6"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-7"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.1"), 0, "keyword-8"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.1"), 1, "keyword-9"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 0, "keyword-10"},
			{time.Duration(60) * time.Second, net.ParseIP("127.0.0.2"), 1, "keyword-11"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-12"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 1, "keyword-13"},
			{time.Duration(120) * time.Second, net.ParseIP("127.0.0.1"), 0, "keyword-14"},
		},
	},
}

func Test_BottomDistribution(t *testing.T) {
	for k, tt := range bottomTests {
		cp := New(generateLists("keyword-", tt.keywordCount), generateLists("127.0.0.", tt.proxyCount), tt.pulse)

		if len(cp) != len(tt.out) {
			t.Errorf("Test %d: Non-equal slice lengths. Got: %d Expected: %d\n", k, len(cp), len(tt.out))
		}
		
		for n:=0; n<len(cp); n++ {
			if cp[n].Time != tt.out[n].Time {
				t.Errorf("Test %d: Time not equal. Got: %d Expected: %d\n", k, cp[n].Time, tt.out[n].Time)
			}
			
			if cp[n].Proxy.String() != tt.out[n].Proxy.String() {
				t.Errorf("Test %d: Proxy not equal. Got: %s Expected: %s\n", k, cp[n].Proxy.String(), tt.out[n].Proxy.String())
			}
			
			if cp[n].Conn != tt.out[n].Conn {
				t.Errorf("Test %d: Connection count not equal. Got: %d Expected: %d\n", k, cp[n].Conn, tt.out[n].Conn)
			}
			
			if cp[n].Conn != tt.out[n].Conn {
				t.Errorf("Test %d: Connection count not equal. Got: %d Expected: %d\n", k, cp[n].Conn, tt.out[n].Conn)
			}
			
			if cp[n].Keyword != tt.out[n].Keyword {
				t.Errorf("Test %d: Keyword not equal. Got: %d Expected: %d\n", k, cp[n].Keyword, tt.out[n].Keyword)
			}			
		}
	}
}

func Test_Distinct(t *testing.T) {
    var tableCr = CrawlPlan{
        {time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
        {time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-1"},
        {time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-2"},
        {time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
        {time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-1"},
        {time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-2"},
        {time.Duration(0) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-0"},
        {time.Duration(60) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-1"},
        {time.Duration(120) * time.Second, net.ParseIP("127.0.0.0"), 0, "keyword-2"},
    }
    
    dist := tableCr.Distinct()
    
    distinctCount := len(dist)
    if distinctCount != 3 {
        t.Errorf("Distinct count wrong. Got: %d Expected: %d\n", distinctCount, 3)
    }

    if dist[0] != 0 {
        t.Errorf("Distinct item 0 is wrong. Got: %d Expected: %d\n", dist[0], 0)
    }

    if dist[1] != 60 {
        t.Errorf("Distinct item 0 is wrong. Got: %d Expected: %d\n", dist[1], 60)
    }

    if dist[2] != 120 {
        t.Errorf("Distinct item 0 is wrong. Got: %d Expected: %d\n", dist[2], 120)
    }
}

func generateLists(prefix string, count int) []string {
	var out []string
	for n:=0; n< count; n++ {
		out = append(out, fmt.Sprintf("%s%d", prefix, n))
	}
	return out
}
