

# PeeVee - Go 1.18+ only (requires Go generics)

![Goreportcard](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat) 
[![Build Status](https://travis-ci.org/migueleliasweb/peevee.svg?branch=master)](https://travis-ci.org/migueleliasweb/peevee) 
[![Coveralls](https://coveralls.io/repos/github/migueleliasweb/peevee/badge.svg?branch=master)](https://coveralls.io/github/migueleliasweb/peevee?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/migueleliasweb/peevee.svg)](https://pkg.go.dev/github.com/migueleliasweb/peevee)

PEEVEE lets you peek into what is happening in real time through Golang channels. It can automatically generate [Prometheus](https://prometheus.io/) metrics about the throughput of channels. Think of it like Unix's *"pv"* (https://linux.die.net/man/1/pv) but in Golang.

## Examples

### Using PeeVee to create both reader and writer channels

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
	pv := peevee.New("myWorkChannel", WithPromMetrics[int]())

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

### Using PeeVee to wrap an existing channel for reading

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
	// this is a channel that gets returned
	// from an SDK or a third party lib, for example
	channelYouDontControl := make(chan bool)

	go func() {
		for {
			// this would be something happening inside the SDK or lib
			channelYouDontControl <- true
			time.Sleep(time.Millisecond * 250)
		}
	}()

	pv := peevee.NewReaderWrap(
		"boolwrap",
		channelYouDontControl, // this time you can pass the existing channel
	)

	// mimic an asynchronous channel reader
	go func() {
		for {
			fmt.Printnl("got:", <-pv.GetReadableChan())
		}
	}()

	runMetricsEndpoint()
	
	// Leave the example running and open "localhost:8080/metrics" on your browser.
	// You will see a metric called "peevee" that is generated in real time
	// when the channel is being read/written.
}
```
