package peevee

import (
	"sync/atomic"
	"time"
	"unsafe"
)

//Config Configs to PeeVee
type Config struct {
	Name         string
	StatsHandler StatsHandler
}

//PeeVee Representation of the PV
type PeeVee struct {
	Name         string
	readChan     chan interface{}
	writeChan    chan interface{}
	statsChan    chan PVStats
	counter      uint64
	counterTime  time.Time
	statsHandler StatsHandler
	messageSize  uintptr
	bitsCounter  uintptr
}

//GetWriteChannel Returns the write channel
func (pv *PeeVee) GetWriteChannel() chan<- interface{} {
	return pv.writeChan
}

//GetReadChannel Returns the read channel
func (pv *PeeVee) GetReadChannel() <-chan interface{} {
	return pv.readChan
}

//procesStats Processes stats and sends them to `pv.statsChan`
func (pv *PeeVee) procesStats(msg interface{}) {
	if pv.messageSize == 0 {
		atomic.StoreUintptr(&pv.messageSize, unsafe.Sizeof(msg))
	}

	atomic.AddUint64(&pv.counter, 1)
	atomic.AddUintptr(&pv.bitsCounter, pv.messageSize)

	if time.Now().After(pv.counterTime.Add(time.Minute)) {
		counter := atomic.LoadUint64(&pv.counter)

		var zeroCounter uint64
		atomic.SwapUint64(&pv.counter, zeroCounter)

		bitsPerSecond := uint64(pv.bitsCounter / 60)
		KbitPerSecond := uint64(bitsPerSecond / 1000)
		MbitPerSecond := uint64(KbitPerSecond / 1000)

		pv.statsChan <- PVStats{
			Name:          pv.Name,
			PerSecond:     uint64(counter / 60),
			BitPerSecond:  bitsPerSecond,
			KbitPerSecond: KbitPerSecond,
			MbitPerSecond: MbitPerSecond,
		}
	}
}

//channelPiper Pipes information between read and write channels
func (pv *PeeVee) channelPiper() {
	for {
		select {
		case msg := <-pv.writeChan:
			pv.procesStats(msg)
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

	if config.StatsHandler == nil {
		pv.statsHandler = NewStdoutStatsHandler()
	} else {
		pv.statsHandler = config.StatsHandler
	}

	go pv.statsHandler.Handle(pv.statsChan)
	go pv.channelPiper()

	return pv
}
