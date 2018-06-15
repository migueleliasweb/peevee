package peevee

type dummyStatsHandler struct{}

func (stats *dummyStatsHandler) Handle(statsChan <-chan PVStats) {
	for {
		<-statsChan
	}
}

type fakeStatsHandlerOkChan struct {
	okChan chan bool
}

func (stats *fakeStatsHandlerOkChan) Handle(statsChan <-chan PVStats) {
	for {
		<-statsChan
		stats.okChan <- true
	}
}
