/*

Usage:
	crawlplan --keywords="./keywords.txt" --proxies="./proxies.txt" --avgJobRuntime=60 --minimumDelay=60 --strategy=bottom --algo="duration" --timePeriod=3600
	crawlplan --keywords="./keywords.txt" --proxies="./proxies.txt" --avgJobRuntime=60 --minimumDelay=60 --strategy=bottom --algo="connections" --maximumConnections=5

*/
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"errors"
	"time"
	"text/tabwriter"

	"stash.stickyeyes.com/groun/crawlrate"
)

var (
	debug 		  *bool = flag.Bool("debug", false, "Switch on debug mode")
	keywordFile   *string = flag.String("keywords", "", "Keyword file")
	proxyFile     *string = flag.String("proxies", "", "proxy file")
	avgJobRuntime *time.Duration = flag.Duration("avgJobRuntime", time.Duration(60) * time.Second, "Average job runtime in seconds. Defaults to 60")
	minimumDelay  *time.Duration = flag.Duration("minimumDelay", time.Duration(0) * time.Second, "A minimum delay between jobs in seconds. Defaults to 0")
	algo          *string = flag.String("algorithm", "duration", "algorithm used [duration|connections]")

	timePeriod *time.Duration = flag.Duration("timePeriod", time.Duration(3600) * time.Second, "The total maximum duration as used by the duration algorithm. Defaults to 3600 seconds (1hr)")
	maximumConnections *int = flag.Int("maximumConnections", 5, "The total maximum connections as used by the connections algorithm. Defaults to ")
)

var (
	badAlgorithm = errors.New("Unrecognised algorithm")
)

func main() {
	flag.Parse()

	if *keywordFile == "" {
		fmt.Println("missing option '-keywords': keyword filename must be specified")
		os.Exit(-1)
	}

	if *proxyFile == "" {
		fmt.Println("missing option '-proxies': proxy filename must be specified")
		os.Exit(-1)
	}

	if *algo == "" {
		fmt.Println("missing option 'algorithm': algorithm must be either 'duration' or 'connections'")
		os.Exit(-1)
	}

	keywords, err := readLines(*keywordFile)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(sort.StringSlice(keywords))	

	proxies, err := readLines(*proxyFile)
	if err != nil {
		log.Fatal(err)
	}

	var p *crawlrate.Pulse
	switch *algo {
	case "duration":
		p, _ = crawlrate.FixedDuration(len(keywords), len(proxies), int(avgJobRuntime.Seconds()), int(minimumDelay.Seconds()), int(timePeriod.Seconds()))
		
	case "connections":
		p, _ = crawlrate.FixedConnections(len(keywords), len(proxies), int(avgJobRuntime.Seconds()), int(minimumDelay.Seconds()), *maximumConnections)
		
	default:
		log.Fatal(badAlgorithm)
	}

	if *debug {
		fmt.Printf("Keywords: %d, %+v\n", len(keywords), keywords)
		fmt.Printf("Proxies: %d, %+v\n", len(proxies), proxies)
		fmt.Printf("Pulse: %+v\n", p)
		fmt.Printf("Number of connections required for each proxy: %d\n", p.Volume)
		fmt.Printf("Maximum runtime per keyword: %ds (%s)\n", int(p.Frequency.Seconds()), p.Frequency.String())
		fmt.Printf("Total duration required to process all keywords: %ds (%s)\n", int(p.Duration.Seconds()), p.Duration.String())
	}

	cp := crawlrate.New(keywords, proxies, p)
	crawlrate.OrderedBy(crawlrate.Start, crawlrate.Proxy, crawlrate.IncreasingConnections).Sort(cp)
	
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, r := range cp {
		fmt.Fprintf(w, "[%ds]\t%s\t%d\t%s\n", int(r.Time.Seconds()), r.Proxy.String(), r.Conn, r.Keyword)
	} 
	w.Flush()
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)

	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))

	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}

		buffer.Write(part)

		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}
