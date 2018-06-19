package peevee

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

//PVStats Per second stats for the Pipe
type PVStats struct {
	Name      string
	PerSecond uint64
}

//StatsHandler Interface for stats handlers
type StatsHandler interface {
	Handle(statsChan <-chan PVStats)
}

//DefaultStatsHandler Default stats handler for PeeVee
type DefaultStatsHandler struct {
	writer io.Writer
}

//Handle Handle the stats
func (stats *DefaultStatsHandler) Handle(statsChan <-chan PVStats) {
	for {
		select {
		case pvStats := <-statsChan:
			jsonBytes, err := json.Marshal(pvStats)

			if err == nil {
				fmt.Fprintln(stats.writer, string(jsonBytes))
			} else {
				log.Panicln(err)
			}
		}
	}
}

//NewStdoutStatsHandler Creates new Stdout handler
func NewStdoutStatsHandler() StatsHandler {
	return &DefaultStatsHandler{
		writer: os.Stdout,
	}
}
