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
	counterMsg   uint64
	counterTime  atomic.Value
	statsHandler StatsHandler
	messageSize  uintptr
	counterBits  uintptr
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

	atomic.AddUint64(&pv.counterMsg, 1)
	atomic.AddUintptr(&pv.counterBits, pv.messageSize)

	//if now is in the future compared to the start time + 30s
	if time.Now().After(pv.counterTime.Load().(time.Time).Add(time.Second * 10)) {
		counterMsg := atomic.LoadUint64(&pv.counterMsg)

		bitsPerSecond := uint64(pv.counterBits / 10)
		KbitPerSecond := uint64(bitsPerSecond / 1000)
		MbitPerSecond := uint64(KbitPerSecond / 1000)

		//reseting internal state
		atomic.StoreUintptr(&pv.counterBits, 0)
		pv.counterTime.Store(time.Now())
		atomic.StoreUint64(&pv.counterMsg, 0)

		pv.statsChan <- PVStats{
			Name:          pv.Name,
			MsgPerSecond:  uint64(counterMsg / 10),
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
		readChan:   make(chan interface{}),
		writeChan:  make(chan interface{}),
		statsChan:  make(chan PVStats, 1),
		counterMsg: uint64(0),
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
