package crawlrate

import (
	"testing"
	"time"
)

var tests = []struct {
	keywordCount, proxyCount, avgJobRuntime, minimumDelay, timePeriod int
	p                                                                 *Pulse
}{
	// NOTE:
	// For ALL TESTS the Pulse Duration should match the timePeriod

	// Volume tests
	{100, 10, 45, 15, 60, &Pulse{10, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
	{100, 10, 45, 15, 120, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},

	// 30 keywords per minute for 4 minutes (last minute is 10 kwpm)
	{100, 10, 45, 15, 240, &Pulse{3, time.Duration(60) * time.Second, time.Duration(240) * time.Second}},

	// 15 kwpm for 3 minutes = 45 keywords: 15 proxies = 1 connection per proxy
	{45, 15, 45, 15, 180, &Pulse{1, time.Duration(60) * time.Second, time.Duration(180) * time.Second}},

	// 15 kwpm for 3 minutes = 45 keywords: 5 proxies = 3 connections per proxy
	{45, 5, 45, 15, 180, &Pulse{3, time.Duration(60) * time.Second, time.Duration(180) * time.Second}},

	// 15 kwpm for 3 minutes = 45 keywords: 2 proxies = 8 connections per proxy
	{45, 2, 45, 15, 180, &Pulse{8, time.Duration(60) * time.Second, time.Duration(180) * time.Second}},

	// 13 kwpm for 8 minutes = 100 keywords: 10 proxies = 2 connections per proxy
	{100, 10, 45, 15, 480, &Pulse{2, time.Duration(60) * time.Second, time.Duration(480) * time.Second}},

	// Frequency expansion tests
	{100, 10, 45, 15, 70, &Pulse{10, time.Duration(70) * time.Second, time.Duration(70) * time.Second}},
	{100, 10, 45, 15, 140, &Pulse{5, time.Duration(70) * time.Second, time.Duration(140) * time.Second}},

	// Algo for testing:
	// 500s/60s = 8m18s (8m with 18s distributed)
	// 100 / 8m = 12.5 rounded up = 13 kwps
	// 13 kwpm for 8m = 100 : 10 proxies = 2 connections per proxy
	{100, 10, 45, 15, 500, &Pulse{2, time.Duration(62) * time.Second, time.Duration(500) * time.Second}},

	{4, 3, 60, 60, 300, &Pulse{1, time.Duration(150) * time.Second, time.Duration(300) * time.Second}},
}

func Test_FixedDuration(t *testing.T) {
	for k, tt := range tests {
		p, err := FixedDuration(tt.keywordCount, tt.proxyCount, tt.avgJobRuntime, tt.minimumDelay, tt.timePeriod)
		if err != nil {
			t.Errorf("Error test: %d Error: %s", k, err)
			continue
		}
		if p.Volume != tt.p.Volume {
			t.Errorf("Error test: %d Volume exp: %d got: %d %v", k, tt.p.Volume, p.Volume, p)
		}
		if p.Frequency.Seconds() != tt.p.Frequency.Seconds() {
			t.Errorf("Error test: %d Frequency exp: %f got: %f %v", k, tt.p.Frequency.Seconds(), p.Frequency.Seconds(), p)
		}
		if p.Duration.Seconds() != tt.p.Duration.Seconds() {
			t.Errorf("Error test: %d Duration exp: %f got: %f %v", k, tt.p.Duration.Seconds(), p.Duration.Seconds(), p)
		}
	}
}

var errorTests = []struct {
	keywordCount, proxyCount, avgJobRuntime, minimumDelay, timePeriod int
	msg                                                               string
}{
	// Duration limit tests
	{100, 10, 45, 15, 50, "total keyword time exceeds duration"},
}

func Test_FixedDurationErrors(t *testing.T) {
	for k, tt := range errorTests {
		_, err := FixedDuration(tt.keywordCount, tt.proxyCount, tt.avgJobRuntime, tt.minimumDelay, tt.timePeriod)
		if err == nil {
			t.Errorf("Error test: %d Expected error", k)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("Error test: %d exp: %s got: %s", k, tt.msg, err.Error())
		}
	}
}
