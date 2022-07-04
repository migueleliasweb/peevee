package peevee

import (
	"strconv"
	"testing"
)

// BenchmarkPeeVeeDefaultOptions [07/2022] 861704|1393 ns/op|0 B/op|0 allocs/op
func BenchmarkPeeVeeDefaultOptions(b *testing.B) {
	pv := New("myintpeevee", WithDefault[int]())

	go func(pv PeeVee[int]) {
		r := pv.GetReadableChan()

		for range r {
		}
	}(pv)

	rw := pv.GetWritableChan()
	for i := 0; i < b.N; i++ {
		rw <- i
	}

}

// BenchmarkPeeVeePromMetrics [07/22] 767596|1806 ns/op|0 B/op|0 allocs/op
func BenchmarkPeeVeePromMetrics(b *testing.B) {
	pv := New("myintpeevee"+strconv.Itoa(b.N), WithPromMetrics[int]())

	go func(pv PeeVee[int]) {
		r := pv.GetReadableChan()

		for range r {
		}
	}(pv)

	rw := pv.GetWritableChan()
	for i := 0; i < b.N; i++ {
		rw <- i
	}

}

// BenchmarkBuiltinChannels [07/22] 2053736|600.3ns/op|0 B/op|0 allocs/op
func BenchmarkBuiltinChannels(b *testing.B) {
	c := make(chan int)

	go func() {
		for range c {
		}
	}()

	for i := 0; i < b.N; i++ {
		c <- i
	}

}
