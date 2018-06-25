package peevee

import (
	"testing"
)

type fakeWriter struct {
	okChan chan bool
}

func (fw *fakeWriter) Write(p []byte) (n int, err error) {
	fw.okChan <- true
	return 0, nil
}

func TestDefaultStatsHandler(t *testing.T) {
	sh := DefaultStatsHandler{
		writer: &fakeWriter{okChan: make(chan bool, 1)},
	}

	statsChan := make(chan PVStats, 1)

	stats := PVStats{
		Name:      "foo",
		PerSecond: uint64(123),
	}

	go sh.Handle(statsChan)
	statsChan <- stats

	//syncronization
	<-sh.writer.(*fakeWriter).okChan

}
