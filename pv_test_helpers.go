package peevee

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
	for {
		select {
		case <-statsChan:
			stats.okChan <- true
			return
		}
	}
}
