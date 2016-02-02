package crawlrate

import (
	"fmt"
	"net"
	"sort"
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

type CrawlRule struct {
	Time    time.Duration
	Proxy   net.IP
	Conn    int
	Keyword string
}

func (cr CrawlRule) String() string {
	return fmt.Sprintf("[%ds]\t%s\t%d\t%s\n", int(cr.Time.Seconds()), cr.Proxy.String(), cr.Conn, cr.Keyword)
}

type lessFunc func(p1, p2 *CrawlRule) bool

// multiSorter implements the Sort interface, sorting the crawlRules within.
type multiSorter struct {
	crawlRules []CrawlRule
	less       []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(crawlRules []CrawlRule) {
	ms.crawlRules = crawlRules
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.crawlRules)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.crawlRules[i], ms.crawlRules[j] = ms.crawlRules[j], ms.crawlRules[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call. We
// could change the functions to return -1, 0, 1 and reduce the
// number of calls for greater efficiency: an exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.crawlRules[i], &ms.crawlRules[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

/*
NOTES:

Benefits of this solution:

- The different scheduling strategies i.e top vs. bottom weighting, left vs. right stacking
    and maximise connections per proxy vs minimise can all be implemented as
    separate functions that can be composed. This introduces a much simpler
    structure, far more flexibility in composition and provides a very clean implementation.

- Uses common stdlib patterns

This is a far better solution; it is also far more performant than anything I've written so far!
*/

func New(keywords, proxies []string, pulse *Pulse) (cr []CrawlRule) {
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

func BottomHeavy(cp []CrawlRule) {
	orderedBy(start, proxy, increasingConnections).Sort(cp)
}

func TopHeavy(cp []CrawlRule) {
	orderedBy(start, proxy, increasingConnections).Sort(cp)
}


// Filter
func Filter(vs []CrawlRule, f func(CrawlRule) bool) []CrawlRule {
    vsf := make([]CrawlRule, 0)
    for _, v := range vs {
        if f(v) {
            vsf = append(vsf, v)
        }
    }
    return vsf
}

func Map(vs []CrawlRule, f func(CrawlRule) CrawlRule) []CrawlRule {
    vsm := make([]CrawlRule, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}
/*
func Distinct(vs []CrawlRule) map[int64]bool {
    var vsm map[int64]bool
    var current time.Duration
    for i, v := range vs {
        if current == nil {
            current = v.Time
        }
        
        if v.Time.Seconds() != current.Seconds()  {
            current = v.Time
            vsm[v.Time.Seconds()] = true
        }
    }
    return vsm
}
*/