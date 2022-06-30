

# PEEVEE Go 1.18+ only (requires Go generics)

[Goreportcard](https://goreportcard.com/badge/github.com/migueleliasweb/peevee) 
[![Build Status](https://travis-ci.org/migueleliasweb/peevee.svg?branch=master)](https://travis-ci.org/migueleliasweb/peevee) 
[![Coveralls](https://coveralls.io/repos/github/migueleliasweb/peevee/badge.svg?branch=master)](https://coveralls.io/github/migueleliasweb/peevee?branch=master)

PEEVEE lets you peek into what is happening in real time throught the Channels in Golang. It can expose [Prometheus](https://prometheus.io/) metrics about the throughput of channels it created. Think of it like Unix's *"pv"* (https://linux.die.net/man/1/pv) in a way.

## Example

```go
package main

import (
	"fmt"
	"time"
	"log"

	"github.com/migueleliasweb/peevee"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func runMetricsEndpoint() {
	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	pv := NewPeeVee("myWorkChannel", WithPromMetrics[int]())

	// mimic an asynchronous channel writer
	go func() {
		for {
			pv.GetWritableChan() <- rand.Int()
			time.Sleep(time.Millisecond * 250)
		}
	}()

	// mimic an asynchronous channel reader
	go func() {
		for {
			fmt.Printnl("int:", <-pv.GetReadableChan())
		}
	}()

	runMetricsEndpoint()
	
	// Leave the example running and open "localhost:8080/metrics" on your browser.
	// You will see a metric called "peevee" that is generated in real time
	// when the channel is being read/written.
}
```
