package peevee

import (
	"testing"
	"time"
)

func TestNewPeeVeeReaderWrap(t *testing.T) {
	// this is a channel that gets returned
	// from an SDK or a third party lib for example
	channelYouDontControl := make(chan bool)

	go func() {
		channelYouDontControl <- true
	}()

	pv := NewReaderWrap(
		"boolwrap",
		channelYouDontControl,
	)

	for {
		select {
		case i := <-pv.GetReadableChan():
			if i != true {
				t.FailNow()
			}

			return
		case <-time.After(time.Second * 5):
			t.Error("timeout reading from channel")
			t.FailNow()
		}
	}
}
