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

func TestNewPeeVeeReaderWrapWithClose(t *testing.T) {
	// this is a channel that gets returned
	// from an SDK or a third party lib for example
	channelYouDontControl := make(chan bool)

	go func() {
		// this would be the third party code sending items
		// to the channel
		channelYouDontControl <- true

		// and this is the channel being closed by the same third party code
		close(channelYouDontControl)
	}()

	pv := NewReaderWrap(
		"boolwrap",
		channelYouDontControl,
	)

	for {
		select {
		case i, more := <-pv.GetReadableChan():
			// now, we purposedly don't return
			// and instead stay in the loop waiting for the channel
			// to be closed
			if !more {
				return
			}

			if i != true {
				t.Error("got wrong item from channel:", i)
				t.FailNow()
			}

		case <-time.After(time.Second * 5):
			t.Error("timeout reading from channel")
			t.FailNow()
		}
	}
}
