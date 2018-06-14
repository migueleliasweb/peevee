package peevee

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

//StatsHandler Interface for stats handlers
type StatsHandler interface {
	Handle(statsChan <-chan PVStats)
}

//DefaultStatsHandler Default stats handler for PeeVee
type DefaultStatsHandler struct{}

//Handle Handle the stats
func (stats *DefaultStatsHandler) Handle(statsChan <-chan PVStats) {
	for {
		select {
		case stats := <-statsChan:
			//quick-n-easy json object
			logMapping := make(map[string]string)

			if stats.Name != "" {
				logMapping["PipeName"] = stats.Name
			}
			logMapping["PerSecond"] = strconv.Itoa(int(stats.PerSecond))
			jsonBytes, err := json.Marshal(logMapping)

			if err == nil {
				fmt.Println(string(jsonBytes))
			} else {
				log.Panicln(err)
			}
		}
	}
}
