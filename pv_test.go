package peevee

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSimpleNewPeeVee(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",
	})

	if pv.Name != "pv1" {
		t.Error("PeeVee name should be pv1")
	}
}

func TestProcessStatsCounter(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",
	})

	pv.procesStats()

	if pv.counter != uint64(1) {
		t.Errorf("Wrong counter value, expecting 1 but got %d", pv.counter)
	}
}

func TestProcessStatsSendToChannel(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",

		//just to make sure we are bypassing the default handler
		StatsHandler: &dummyStatsHandler{},
	})

	pv.counterTime = time.Now().Add(time.Minute * -2)

	pv.procesStats()

	stats := <-pv.GetStatsChannel()

	if stats.Name != "pv1" {
		t.Errorf("Wrong stats.Name value, expecting 'pv1' but got %s", stats.Name)
	}

	if stats.PerSecond != uint64(0) {
		t.Errorf("Wrong stats.PerSecond value, expecting 0 but got %d", stats.PerSecond)
	}
}

func TestProcessStatsReceivesFromChannel(t *testing.T) {
	config := Config{
		Name:         "pv1",
		StatsHandler: &fakeStatsHandlerOkChan{},
	}

	pv := NewPeeVee(config)

	pv.counterTime = time.Now().Add(time.Minute * -2)
	pv.procesStats()

	timeoutChan := time.After(time.Second)

	select {
	case <-config.StatsHandler.(*fakeStatsHandlerOkChan).okChan:
		return
	case <-timeoutChan:
		t.Error("Waited too long for the channel to return")
	}
}

func TestChannelPiper(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",

		//just to make sure we are bypassing the default handler
		StatsHandler: func(statsChan <-chan PVStats) {

			counter := 0

			// just drain all stats asap to unblock procesStats
			for {
				select {
				case <-statsChan:
					counter++
				}
			}
		},
	})

	//hack into time
	pv.counterTime = time.Now().Add(time.Minute * -2)

	readCounter := uint64(0)
	total := uint64(1234)
	okChan := make(chan bool, 1)

	go func() {

		for {
			select {
			case <-pv.GetReadChannel():
				atomic.AddUint64(&readCounter, 1)
				if readCounter == total {
					okChan <- true
					return
				}
			}
		}
	}()

	for index := uint64(0); index < total; index++ {
		pv.writeChan <- true
	}

	<-okChan

	//mildly repeated but still...
	if readCounter != total {
		t.Errorf("Wrong total number of reads expected %d got %d", total, readCounter)
	}

}
