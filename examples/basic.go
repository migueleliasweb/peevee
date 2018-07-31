package main

import (
	"log"
	"time"

	"github.com/migueleliasweb/peevee"
)

func basic() {
	config := peevee.Config{
		Name:           "my-string-channel",
		StatsFrequency: uint64(5),
	}

	pv := peevee.NewPeeVee(config)

	go func() {
		for {
			select {
			case msg := <-pv.GetReadChannel():
				log.Println(msg)
			}
		}
	}()

	log.Println("Printing channel stats every", config.StatsFrequency, "seconds...")

	for {
		pv.GetWriteChannel() <- "PeeVee is AWESOME"
		time.Sleep(time.Millisecond * 200)
	}
}
