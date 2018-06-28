package basic

import (
	"fmt"
	"time"

	"github.com/migueleliasweb/peevee"
)

// WithoutPeeVee Basic example 1 without PeeVee
func WithoutPeeVee() {
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

// WithPeeVee Basic example 1 with PeeVee
func WithPeeVee() {
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
