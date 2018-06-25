package peevee

import (
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
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

	pv.procesStats(true)

	if pv.counter != uint64(1) {
		t.Errorf("Wrong counter value, expecting 1 but got %d", pv.counter)
	}
}

func TestProcessStatsMessageSize(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",
	})

	pv.procesStats(true)
	expected := unsafe.Sizeof(true)

	if pv.messageSize != expected {
		t.Errorf(
			"Wrong message size value, expecting %d but got %d",
			expected,
			pv.messageSize,
		)
	}
}

func TestProcessStatChannel(t *testing.T) {
	handler := okChanStatsHandler{
		okChan: make(chan bool),
	}

	pv := NewPeeVee(Config{
		Name: "pv1",

		//just to make sure we are bypassing the default handler
		StatsHandler: &handler,
	})

	//we need to internally change the time as the method checks if
	//it needs to send a msg to the channel
	pv.counterTime = time.Now().Add(time.Minute * -2)
	pv.procesStats(true)

	for {
		select {
		case <-pv.statsHandler.(*okChanStatsHandler).okChan:
			return
		case <-time.After(time.Second * 5):
			t.Error("Stats handler took too long to process")
		}
	}
}

func TestChannelPiper(t *testing.T) {
	pv := NewPeeVee(Config{
		Name: "pv1",

		//just to make sure we are bypassing the default handler
		StatsHandler: &dummyStatsHandler{},
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
