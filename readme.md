# CrawlPlan

An experiment in calculating fixed schedules for a Google crawler.

## Overview

Given a list of keywords and a list of proxies, calculate a schedule constrained
to either a fixed number of connections per proxy or a fixed duration to crawl
all the keywords.

Each algorithm returns a Pulse. A Pulse consists of the following properties:
 - Volume:    The number of connections required for each proxy
 - Frequency: The maximum runtime per keyword
 - Duration:  The total duration required to process all keywords

From a Pulse, a crawl plan can be created.

### Fixed Duration

FixedDuration algorithm favours a fixed duration over a variable connection 
count per proxy. This algorithm will tell you how many connections each proxy 
needs in order to complete the keywords with the proxies provided within the 
fixed time duration.

### Fixed Connection Count

FixedConnections algorithm favours fixed connection count per proxy over a 
variable duration. This algorithm will tell you the total time duration it will 
take to complete the keywords given the proxies and limiting them to a fixed 
number of connections. If the user supplies a connectionCount greater than that 
required, then the returned pulse will contain a corrected connectionCount 
(volume).

## Crawl Plan

A crawl plan can be either be top heavy or bottom heavy. A bottom heavy crawl 
plan will plan the maxium pulse crawls at the bottom of the time period i.e. as 
soon as possible, where as a top heavy crawl plan will favor the top of the 
period.

### Bottom Heavy
```
  0s * * * * *
 60s * * * * *
120s * * *
```

### Top Heavy
```
  0s * * *
 60s * * * * *
120s * * * * *
```

Each tick is a proxy and a connection count e.g.

```
  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 2 2 2      ->  6 kwpm
```

Points of interest are how we distribute the final minute. For example, compare 
the following two crawl plans; one reduces the number of connections across all 
proxies, whereas the other will reprieve 2 of the proxies from any connections:

```
  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 2 2 1      ->  5 kwpm
```

```
  0s 2 2 2 2 2  -> 10 kwpm
 60s 2 2 2 2 2  -> 10 kwpm
120s 1 1 1 1 1  ->  5 kwpm
```

### Crawl Plan calculations

A crawl plan is created from a Pulse. The following calculations are made:

 - numberOfRows = pulse.Duration / pulse.Frequency
 - numberOfColumns = proxies []string  cellValue = pulse.Volume
 - keywordCount (to calculate last row)

which means:

 - keyword count =/= cellValue * numberOfColumns * numberOfRows

