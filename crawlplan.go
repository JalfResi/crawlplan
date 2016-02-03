package crawlrate

import (
	"net"
	"time"
)

/*
The crawl plan
==============

A crawl plan can be either be top heavy or bottom heavy. A bottom heavy crawl plan will
plan the maxium pulse crawls at the bottom of the time period i.e. as soon as possible,
where as a top heavy crawl plan will favor the top of the hour.

e.g.

bottom heavy:
  0s * * * * *
 60s * * * * *
120s * * *

top heavy:
  0s * * *
 60s * * * * *
120s * * * * *

Each tick is a proxy and a a connection count e.g.

  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 2 2 2      ->  6 kwpm

Points of interest are how we distribute the final minute. For example, compare the following
two crawl plans; one reduces the number of connections across all proxies, whereas the other
will reprieve 2 of the proxies from any connections:

  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 2 2 1      ->  5 kwpm

  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 1 1 1 1 1  ->  5 kwpm

-----------------

A crawl plan is created from a Pulse. The following calculations are made:

    numberOfRows = pulse.Duration / pulse.Frequency

    numberOfColumns = proxies []string  cellValue = pulse.Volume

    keywordCount (to calculate last row)

which means:

    keyword count =/= cellValue * numberOfColumns * numberOfRows



*/


type CrawlPlan []CrawlRule

func New(keywords, proxies []string, pulse *Pulse) (cr CrawlPlan) {
	/*
		NOTE:
		May have to add a 'last row' check so that I can apply different distribution
		algorithms e.g distribute remaining keywords across all proxies vs. cluster as
		many connections to as few proxies as possible. Bear in mind that some tests only
		have a single row, so there is no "last row" to balance. What do we do in those
		circumstances?
	*/
	var currentKeyword int = 0

outerLoop:
	for t := 0; t < int(pulse.Duration.Seconds()); t = t + int(pulse.Frequency.Seconds()) {
		for _, proxy := range proxies {
			for conn := 0; conn < pulse.Volume; conn++ {
				cr = append(cr, CrawlRule{time.Duration(t) * time.Second, net.ParseIP(proxy), conn, keywords[currentKeyword]})
				currentKeyword++
				if currentKeyword >= len(keywords) {
					break outerLoop
				}
			}
		}
	}

	orderedBy(start, proxy, increasingConnections).Sort(cr)
	return
}

// Filter
func (cp CrawlPlan) Filter(t float64, f func(float64, CrawlRule) bool) CrawlPlan {
    vsf := make(CrawlPlan, 0)
    for _, v := range cp {
        if f(t, v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

// Map
func (cp CrawlPlan) Map(f func(CrawlRule) CrawlRule) {
    for _, v := range cp {
        v = f(v)
    }
}


func (cp CrawlPlan) Distinct() map[float64]bool {
    vsm := make(map[float64]bool)
    for _, v := range cp {
        vsm[v.Time.Seconds()] = true
    }
    return vsm
}


func start(c1, c2 *CrawlRule) bool {
	return c1.Time.Seconds() < c2.Time.Seconds()
}

func proxy(c1, c2 *CrawlRule) bool {
	return c1.Proxy.String() < c2.Proxy.String()
}

func increasingConnections(c1, c2 *CrawlRule) bool {
	return c1.Conn < c2.Conn
}

func decreasingConnections(c1, c2 *CrawlRule) bool {
	return c1.Conn > c2.Conn // Note: > orders downwards.
}

func BottomHeavy(cp CrawlPlan) {
	orderedBy(start, proxy, increasingConnections).Sort(cp)
}

func TopHeavy(cp CrawlPlan) {
	orderedBy(start, proxy, increasingConnections).Sort(cp)
}
