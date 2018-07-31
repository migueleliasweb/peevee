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
