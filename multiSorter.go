package crawlrate

import (
	"sort"
)

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

type lessFunc func(p1, p2 *CrawlRule) bool

// multiSorter implements the Sort interface, sorting the crawlRules within.
type multiSorter struct {
	crawlRules CrawlPlan
	less       []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(crawlRules CrawlPlan) {
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
