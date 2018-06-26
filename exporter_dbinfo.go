package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func init() {
	RegisterExporter("dbinfo", newExporterDbinfo)
}

type DbinfoItems struct {
	Name           string `json:"name"`
	StorageEnabled bool   `json:"storage_enabled"`
}

type DbinfoResponse struct {
	Items      []DbinfoItems `json:"items"`
	TotalItems int           `json:"total_items"`
}

type exporterDbinfo struct {
	dbinfoMetricsGauge map[string]*prometheus.GaugeVec
}

var (
	dbinfoLabels   = []string{"db"}
	dbinfoGaugeVec = map[string]*prometheus.GaugeVec{
		"total_namespaces": newGaugeVec("total_namespaces", "total namespaces count in the database", dbinfoLabels),
	}
)

func newExporterDbinfo() Exporter {
	return exporterDbinfo{
		dbinfoMetricsGauge: dbinfoGaugeVec,
	}
}

func (e exporterDbinfo) String() string {
	return "Exporter dbinfo"
}

func (e exporterDbinfo) Collect(ch chan<- prometheus.Metric) error {
	dbinfoData, _ := apiGetNamespacesList(config)
	dbinfo := DbinfoResponse{}
	if err := json.Unmarshal(dbinfoData, &dbinfo); err != nil {
		return err
	}

	for gaugeKey, gauge := range e.dbinfoMetricsGauge {
		if strings.Compare(gaugeKey, "total_namespaces") == 0 {
			gauge.WithLabelValues(config.ReindexerDBName).Set(float64(dbinfo.TotalItems))
		}

		gauge.Collect(ch)
	}

	return nil
}

func (e exporterDbinfo) Describe(ch chan<- *prometheus.Desc) {
	for _, dbinfoMetric := range e.dbinfoMetricsGauge {
		dbinfoMetric.Describe(ch)
	}
}
