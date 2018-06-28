

# PEEVEE ![Goreportcard](https://goreportcard.com/badge/github.com/migueleliasweb/peevee) [![Build Status](https://travis-ci.org/migueleliasweb/peevee.svg?branch=master)](https://travis-ci.org/migueleliasweb/peevee) [![Coveralls](https://coveralls.io/repos/github/migueleliasweb/peevee/badge.svg?branch=master)](https://coveralls.io/github/migueleliasweb/peevee?branch=master)

PEEVEE allows you to peek into what is happening in real time throught the Channels in Golang. Think like Unix *"pv"* in bash (https://linux.die.net/man/1/pv).

### Basic example: What if we would like to know how many messages are passing throught the 'queue' channel?

```go
package main

import (
	"fmt"
	"time"

	"github.com/migueleliasweb/peevee"
)

func withoutPeevee() {
	queue := make(chan string)

	go func() {
		for {
			select {
			case msg := <-queue:
				fmt.Println(msg)
			}
		}
	}()

	for {
		queue <- "PEEVEE is AWESOME"
		time.Sleep(time.Millisecond * 200)
	}
}

func withPeevee() {
	pv := peevee.NewPeeVee(peevee.Config{Name: "my-string-channel"})

	go func() {
		for {
			select {
			case msg := <-pv.GetReadChannel():
				fmt.Println(msg)
			}
		}
	}()

	for {
		pv.GetWriteChannel() <- "PEEVEE is AWESOME"
		time.Sleep(time.Millisecond * 200)
	}
}

func main() {
	withoutPeevee()
	//withPeevee()
}

```
