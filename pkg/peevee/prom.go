package peevee

import (
	"github.com/prometheus/client_golang/prometheus"
)

var metricName = "peevee"

type WithPromMetricsConfig struct {
	// UseRegisterer If not set will default to `prometheus.DefaultRegisterer`
	UseRegisterer prometheus.Registerer
}

// SetMetricName Use this function to set the metric name.
//
// Please do not use this function after calling `NewPeeVee` as it will
// generate inconsistent metric names.
//
// Only call this function once.
func SetMetricName(name string) {
	metricName = name
}

func WithPromMetrics[T any](c ...WithPromMetricsConfig) PVOptions[T] {
	return func(pv *PeeVee[T]) {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: metricName,
			ConstLabels: prometheus.Labels{
				"channel": pv.Name,
			},
		})

		registerer := prometheus.DefaultRegisterer

		if len(c) == 1 {
			if c[0].UseRegisterer != nil {
				registerer = c[0].UseRegisterer
			}
		}

		registerer.MustRegister(counter)

		ccp := CallbackChannelPiper[T]{
			ReadChan:  pv.readChan,
			WriteChan: pv.writeChan,
			callback: func(t T) {
				counter.Inc()
			},
		}

		go ccp.Pipe()
	}
}
