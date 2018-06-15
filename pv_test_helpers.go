package peevee

import (
	"fmt"
	"sync/atomic"
)

type dummyStatsHandler struct{}

func (stats *dummyStatsHandler) Handle(statsChan <-chan PVStats) {
	for {
		fmt.Println("aaa")
		<-statsChan
		fmt.Println("bbb")
	}
}

type fakeStatsHandlerOkChan struct {
	okChan chan bool
}

func (stats *fakeStatsHandlerOkChan) Handle(statsChan <-chan PVStats) {
	stats.okChan = make(chan bool, 1)
	c := uint64(0)

	for {
		select {
		case <-statsChan:
			atomic.AddUint64(&c, 1)
			stats.okChan <- true
		}
	}
}
