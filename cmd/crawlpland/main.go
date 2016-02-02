package main
/*
crawlpland

RPC Server for merging multiple crawlplans. Provides several methods for creating 
mutiple plans and producing a merged crawl plan.

Requires a client (crawlplanc) to make the method calls.

RPC Methods
- AddDuration(avgJobRuntime, minimumDelay, keywords, proxies, timePeriod)
- AddConnections(avgJobRuntime, minimumDelay, keywords, proxies, maximumConnections)
- CreateCrawlPlan()

*/

import(
    "net/rpc"
    "net"
    "net/http"
    "log"
    "./rpc"
)

s := &Server{
    sessions: make([SessionId]*Session),
}

rpc.Register(s)
rpc.HandleHTTP()
l, e := net.Listen("tcp", ":1234")
if e != nil {
	log.Fatal("listen error:", e)
}
go http.Serve(l, nil)
