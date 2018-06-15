package peevee

import (
	"sync/atomic"
	"time"
)

//Config Configs to PeeVee
type Config struct {
	Name         string
	StatsHandler StatsHandler
}

//PeeVee Representation of the PV
type PeeVee struct {
	Name        string
	readChan    chan interface{}
	writeChan   chan interface{}
	statsChan   chan PVStats
	counter     uint64
	counterTime time.Time
}

//GetWriteChannel Returns the write channel
func (pv *PeeVee) GetWriteChannel() chan<- interface{} {
	return pv.writeChan
}

//GetReadChannel Returns the read channel
func (pv *PeeVee) GetReadChannel() <-chan interface{} {
	return pv.readChan
}

//GetStatsChannel Returns the stats channel
func (pv *PeeVee) GetStatsChannel() chan PVStats {
	return pv.statsChan
}

//procesStats Processes stats and sends them to `pv.statsChan`
func (pv *PeeVee) procesStats() {
	atomic.AddUint64(&pv.counter, 1)

	if time.Now().After(pv.counterTime.Add(time.Minute)) {
		counter := atomic.LoadUint64(&pv.counter)

		var zeroCounter uint64
		atomic.SwapUint64(&pv.counter, zeroCounter)

		pv.statsChan <- PVStats{
			Name:      pv.Name,
			PerSecond: uint64(counter / 60),
		}
	}
}

//channelPiper Pipes information between read and write channels
func (pv *PeeVee) channelPiper() {
	for {
		select {
		case msg := <-pv.writeChan:
			pv.procesStats()
			pv.readChan <- msg
		}
	}
}

//NewPeeVee Configures and returns a new PeeVee
func NewPeeVee(config Config) PeeVee {
	pv := PeeVee{
		readChan:    make(chan interface{}),
		writeChan:   make(chan interface{}),
		statsChan:   make(chan PVStats, 1),
		counterTime: time.Now(),
		counter:     uint64(0),
	}

	if config.Name != "" {
		pv.Name = config.Name
	}

	var statsHandler StatsHandler

	if config.StatsHandler == nil {
		statsHandler = NewStdoutStatsHandler()
	} else {
		statsHandler = config.StatsHandler
	}

	go statsHandler.Handle(pv.GetStatsChannel())
	go pv.channelPiper()

	return pv
}
