package main

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	exportersMu       sync.RWMutex
	exporterFactories = make(map[string]func() Exporter)
)

//RegisterExporter makes an exporter available by the provided name.
func RegisterExporter(name string, f func() Exporter) {
	exportersMu.Lock()
	defer exportersMu.Unlock()

	if f == nil {
		panic("exporterFactory is nil")
	}

	exporterFactories[name] = f
}

type exporter struct {
	mutex    sync.RWMutex
	upMetric prometheus.Gauge
	exporter []Exporter
}

//Exporter interface for prometheus metrics. Collect is fetching the data and therefore can return an error
type Exporter interface {
	Collect(ch chan<- prometheus.Metric) error
	Describe(ch chan<- *prometheus.Desc)
}

func newExporter() *exporter {
	enabledExporter := []Exporter{}
	for _, e := range config.EnabledExporters {
		enabledExporter = append(enabledExporter, exporterFactories[e]())
	}

	return &exporter{
		upMetric: newGauge("up", "The last scrape of Reindexer is successful."),
		exporter: enabledExporter,
	}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, ex := range e.exporter {
		ex.Describe(ch)
	}

	e.upMetric.Describe(ch)
	BuildInfo.Describe(ch)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	// Protect metrics from concurrent collects.
	e.mutex.Lock()
	defer e.mutex.Unlock()

	start := time.Now()
	allUp := true

	for _, ex := range e.exporter {
		err := ex.Collect(ch)
		if err != nil {
			allUp = false
		}
	}

	BuildInfo.Collect(ch)

	if allUp {
		e.upMetric.Set(1)
	} else {
		e.upMetric.Set(0)
	}
	e.upMetric.Collect(ch)
	log.WithField("duration", time.Since(start)).Info("Metrics updated")

}
