package crawlrate

import (
	"errors"
	"math"
	"time"
)

// FixedDuration algorithm favours a fixed duration over a variable connection count per proxy.
// This algorithm will tell you how many connections each proxy needs in order to complete
// the keywords with the proxies provided within the fixed time duration.
func FixedDuration(keywordCount, proxyCount, avgJobRuntime, minimumDelay, timePeriod int) (*Pulse, error) {

	// Divide the number of keywords between the number of proxies and round
	// up. This gives us our keywordCountPerProxy.
	keywordCountPerProxy := keywordCount / proxyCount

	// Calculate the total time a keyword will take, including its delay.
	totalKeywordTime := avgJobRuntime + minimumDelay

	if totalKeywordTime > timePeriod {
		return nil, errors.New("total keyword time exceeds duration")
	}

	// Determine the number of times we can fit the totalKeywordTime into our
	// time range. This gives us the total number of keywords we can fit into
	// a channel for our timeRange. We round this figure down.
	keywordsPerChannel := timePeriod / totalKeywordTime

	// We now need to work out how much of our timeRange is left over. We can
	// later use this to bump up the minimumDelay of each keyword to fill the
	// remaining space.
	timeRangeRemainder := timePeriod - (totalKeywordTime * keywordsPerChannel)

	// Distribute the timeRangeRemainder over the minimumDelay. We round down
	// here to ensure that the total combined totalKeywordTime is still less
	// than our timeRange.
	minimumDelayIncrement := timeRangeRemainder / keywordsPerChannel

	// Increment the totalKeywordTime to include our delay adjustment. This
	// gives us the frequency of ticks.
	tickFrequency := totalKeywordTime + minimumDelayIncrement

	// Take the number of keywords we can fit into a single channel
	// (keywordsPerChannel) and the total number of keywords. Round up to the
	// nearest whole number. This is the number of channels we need, and
	// consequently, the number of keywords per tick we need to process per
	// tick (volume).
	return &Pulse{
		Volume:    int(math.Ceil(float64(float64(keywordCountPerProxy) / float64(keywordsPerChannel)))),
		Frequency: time.Duration(tickFrequency) * time.Second,
		Duration:  time.Duration(timePeriod) * time.Second,
	}, nil
}
