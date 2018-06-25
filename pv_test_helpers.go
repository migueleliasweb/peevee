package peevee

import "fmt"

type dummyStatsHandler struct{}

func (stats *dummyStatsHandler) Handle(statsChan <-chan PVStats) {
	for {
		<-statsChan
	}
}

type okChanStatsHandler struct {
	okChan chan bool
}

func (stats *okChanStatsHandler) Handle(statsChan <-chan PVStats) {
	fmt.Println("Started handle")

	for {
		select {
		case <-statsChan:
			fmt.Println("okChan")
			stats.okChan <- true
			return
		}
	}
}
