package crawlrate

import (
	"math"
	"time"
)

// FixedConnections algorithm favours fixed connection count per proxy over a variable duration.
// This algorithm will tell you the total time duration it will take to complete the keywords
// given the proxies and limiting them to a fixed number of connections.
// If the user supplies a connectionCount greater than that required, then the returned pulse
// will contain a corrected connectionCount (volume).
func FixedConnections(keywordCount, proxyCount, avgJobRuntime, minimumDelay, connectionCount int) (*Pulse, error) {

	// Divide the number of keywords between the number of proxies and round
	// up. This gives us our keywordCountPerProxy.
	keywordCountPerProxy := int(keywordCount / proxyCount)

	// Calculate the total number of keywords per channel
	keywordsPerChannel := int(math.Ceil(float64(keywordCountPerProxy) / float64(connectionCount)))

	/*
		x := float64(keywordCountPerProxy) / float64(connectionCount)
		fmt.Printf("x: %f keywordsPerChannel: %d\n", x, int(math.Ceil(x)))
	*/

	// Calculate the total time a keyword will take, including its delay.
	totalKeywordTime := avgJobRuntime + minimumDelay

	// Determine the total time period required to process all keywords
	timePeriod := keywordsPerChannel * totalKeywordTime

	return &Pulse{
		Volume:    int(math.Ceil(float64(keywordCountPerProxy) / float64(keywordsPerChannel))),
		Frequency: time.Duration(totalKeywordTime) * time.Second,
		Duration:  time.Duration(timePeriod) * time.Second,
	}, nil
}
