package peevee

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestWithPromMetricsSimple(t *testing.T) {
	NewPeeVee("myintpeevee", WithPromMetrics[int]())
}

func TestWithPromMetrics(t *testing.T) {
	r := prometheus.NewRegistry()

	pv := NewPeeVee("myintpeevee", WithPromMetrics[int](WithPromMetricsConfig{
		UseRegisterer: r,
	}))

	go func() {
		pv.GetWritableChan() <- 1234
	}()

	func() {
		for {
			select {
			case i := <-pv.GetReadableChan():
				if i != 1234 {
					t.FailNow()
				}

				return
			case <-time.After(time.Second * 5):
				t.Error("timeout reading from channel")
				t.FailNow()
			}
		}
	}()

	// at this point the msg was read, and hopefully, the counter incremented... let's check

	metricFamily, err := r.Gather()

	if err != nil {
		t.Error("failed gathering default registry:", err)
		t.FailNow()
	}

	if *(metricFamily[0].Name) != metricName {
		t.Errorf("got wrong metric name: want %s, got %s", metricName, *(metricFamily[0].Name))
	}

	if *(metricFamily[0].Metric[0].Label[0].Value) != "myintpeevee" {
		t.Errorf("got wrong label value for the channel label: got %s, want myintpeevee", *(metricFamily[0].Metric[0].Label[0].Value))
	}

	if *(metricFamily[0].Metric[0].Counter.Value) != float64(1) {
		t.Errorf("got wrong metric value: want 1, got %f", *(metricFamily[0].Metric[0].Counter.Value))
	}
}
