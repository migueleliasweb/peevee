# ***PeeVee*** [![Goreportcard](https://goreportcard.com/badge/github.com/migueleliasweb/peevee)](https://goreportcard.com/report/github.com/migueleliasweb/peevee) [![Build Status](https://travis-ci.org/migueleliasweb/peevee.svg?branch=master)](https://travis-ci.org/migueleliasweb/peevee) [![Coveralls](https://coveralls.io/repos/github/migueleliasweb/peevee/badge.svg?branch=master)](https://coveralls.io/github/migueleliasweb/peevee?branch=master)

PEEVEE allows you to peek into what is happening in real time throught the Channels in Golang. Think it works like Unix **pv** in bash (https://linux.die.net/man/1/pv) **BUT IN Golang**.

### Basic example

```go
package main

import (
	"log"
	"time"

	"github.com/migueleliasweb/peevee"
)

func main() {
    pv := peevee.NewPeeVee(peevee.Config{Name: "my-string-channel"})

	go func() {
		for {
			select {
			case msg := <-pv.GetReadChannel():
				log.Println(msg)
			}
		}
	}()

    log.Println("Printing channel stats every 10 seconds...")

	for {
		pv.GetWriteChannel() <- "PeeVee is AWESOME"
		time.Sleep(time.Millisecond * 200)
	}
}
```
