package peevee

import (
	"testing"
	"time"
)

//END TO END TEST
//we must not use internal values or method in here
//this simulates a high level usage of the package
func TestE2E(t *testing.T) {

	pv := NewPeeVee(Config{
		Name: "E2E",
	})

	go func() {
		pv.GetWriteChannel() <- true
	}()

	for {
		select {
		case <-pv.GetReadChannel():
			return
		case <-time.After(time.Second * 4):
			t.Error("GetReadChannel took too long to receive the message")
		}
	}
}
