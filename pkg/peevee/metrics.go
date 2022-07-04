package peevee

// import (
// 	"io"
// 	"sync/atomic"
// 	"time"
// )

// type WithMetricsWriterConfig struct {
// 	Writer     io.Writer
// 	FlushTimer time.Duration
// }

// type ThroughputCollector struct {
// 	PVName     string
// 	Counter    *int64
// 	FlushTimer time.Duration
// }

// func (tc *ThroughputCollector) Increment() {
// 	atomic.AddInt64(tc.Counter, 1)
// }

// func (tc *ThroughputCollector) Run() {
// 	for {
// 		<-time.After(tc.FlushTimer)

// 		// reset timer
// 		// output metric to writer
// 	}
// }

// func WithMetricsWriter[T any](c WithMetricsWriterConfig) PVOptions[T] {
// 	return func(pv *PeeVee[T]) {
// 		tc := ThroughputCollector{
// 			PVName: pv.Name,
// 		}

// 		ccp := CallbackChannelPiper[T]{
// 			ReadChan:  pv.readChan,
// 			WriteChan: pv.writeChan,
// 			callback: func(t T) {
// 				tc.Increment()
// 				c.Writer.Write([]byte("callback called"))
// 			},
// 		}

// 		go tc.Run()
// 		go ccp.Pipe()
// 	}
// }
