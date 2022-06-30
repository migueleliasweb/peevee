package peevee

import (
	"testing"
	"time"
)

func TestCallbackChannelPiper(t *testing.T) {
	pv := NewPeeVee("myintpeevee", WithCallback(func(i int) {}))

	go func() {
		pv.GetWritableChan() <- 1234
	}()

	for {
		select {
		case i := <-pv.GetReadableChan():
			if i != 1234 {
				t.FailNow()
			}

			return
		case <-time.After(time.Second * 5):
			t.Error("timeout reading from channel")
			t.FailNow()
		}
	}
}

func TestDefaultChannelPiper(t *testing.T) {
	pv := NewPeeVee[int]("myintpeevee")

	go func() {
		pv.GetWritableChan() <- 1234
	}()

	for {
		select {
		case i := <-pv.GetReadableChan():
			if i != 1234 {
				t.FailNow()
			}

			return
		case <-time.After(time.Second * 5):
			t.Error("timeout reading from channel")
			t.FailNow()
		}
	}
}
