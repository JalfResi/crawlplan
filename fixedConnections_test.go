package crawlrate

import (
	"testing"
	"time"
)

var fcTests = []struct {
	keywordCount, proxyCount, avgJobRuntime, minimumDelay, connectionCount int
	p                                                                      *Pulse
}{
	{10, 10, 45, 15, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
	{10, 5, 45, 15, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{10, 5, 45, 15, 2, &Pulse{2, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},

	// Even if the user supplies a larger connection count than is required, the pulse returned
	// should return the actual (less than specified) connection count.

	{10, 10, 45, 15, 2, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
	{10, 10, 45, 15, 10, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
	{10, 10, 45, 15, 100000, &Pulse{1, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},

	// Duration tests
	{100, 10, 45, 15, 10, &Pulse{10, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
	{100, 10, 45, 15, 9, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{100, 10, 45, 15, 8, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{100, 10, 45, 15, 7, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{100, 10, 45, 15, 6, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{100, 10, 45, 15, 5, &Pulse{5, time.Duration(60) * time.Second, time.Duration(120) * time.Second}},
	{100, 10, 45, 15, 4, &Pulse{4, time.Duration(60) * time.Second, time.Duration(180) * time.Second}},
	{100, 10, 45, 15, 3, &Pulse{3, time.Duration(60) * time.Second, time.Duration(240) * time.Second}},
	{100, 10, 45, 15, 2, &Pulse{2, time.Duration(60) * time.Second, time.Duration(300) * time.Second}},
	{100, 10, 45, 15, 1, &Pulse{1, time.Duration(60) * time.Second, time.Duration(600) * time.Second}},

	// Frequency expansion tests
	// keywordCount, proxyCount, avgJobRuntime, minimumDelay, connectionCount
	// Volume, Frequency, Duration
	// {100, 10, 45, 15, 10, &Pulse{10, time.Duration(60) * time.Second, time.Duration(60) * time.Second}},
}

func Test_FixedConnections(t *testing.T) {
	for k, tt := range fcTests {
		p, err := FixedConnections(tt.keywordCount, tt.proxyCount, tt.avgJobRuntime, tt.minimumDelay, tt.connectionCount)
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
