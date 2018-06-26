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
	counterTime  atomic.Value
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

	if time.Now().After(pv.counterTime.Load().(time.Time).Add(time.Second * 30)) {
		counter := atomic.LoadUint64(&pv.counter)

		var zeroCounter uint64
		atomic.SwapUint64(&pv.counter, zeroCounter)

		bitsPerSecond := uint64(pv.bitsCounter / 30)
		KbitPerSecond := uint64(bitsPerSecond / 1000)
		MbitPerSecond := uint64(KbitPerSecond / 1000)

		//reseting internal state
		atomic.StoreUintptr(&pv.bitsCounter, 0)
		pv.counterTime.Store(time.Now())

		pv.statsChan <- PVStats{
			Name:          pv.Name,
			MsgPerSecond:  uint64(counter / 60),
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
		readChan:  make(chan interface{}),
		writeChan: make(chan interface{}),
		statsChan: make(chan PVStats, 1),
		counter:   uint64(0),
	}

	pv.counterTime.Store(time.Now())

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
