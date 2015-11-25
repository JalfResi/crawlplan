package crawlrate

import (
	"time"
)

type Pulse struct {
	Volume    int           // Number of connections required for each proxy
	Frequency time.Duration // Maximum runtime per keyword
	Duration  time.Duration // Total duration required to process all keywords
}
