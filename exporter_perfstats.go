package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func init() {
	RegisterExporter("perfstats", newExporterPerfstats)
}

type PerfstatsResponse struct {
	Items []struct {
		Name    string `json:"name"`
		Updates struct {
			TotalQueriesCount    int `json:"total_queries_count"`
			TotalAvgLatencyUs    int `json:"total_avg_latency_us"`
			TotalAvgLockTimeUs   int `json:"total_avg_lock_time_us"`
			LastSecQPS           int `json:"last_sec_qps"`
			LastSecAvgLockTimeUs int `json:"last_sec_avg_lock_time_us"`
			LastSecAvgLatencyUs  int `json:"last_sec_avg_latency_us"`
		} `json:"updates"`
		Selects struct {
			TotalQueriesCount    int `json:"total_queries_count"`
			TotalAvgLatencyUs    int `json:"total_avg_latency_us"`
			TotalAvgLockTimeUs   int `json:"total_avg_lock_time_us"`
			LastSecQPS           int `json:"last_sec_qps"`
			LastSecAvgLockTimeUs int `json:"last_sec_avg_lock_time_us"`
			LastSecAvgLatencyUs  int `json:"last_sec_avg_latency_us"`
		} `json:"selects"`
	} `json:"items"`
	TotalItems int `json:"total_items"`
}

type exporterPerfstats struct {
	perfstatsMetricsGauge map[string]*prometheus.GaugeVec
}

var (
	perfstatsLabels   = []string{"namespace", "action"}
	perfstatsGaugeVec = map[string]*prometheus.GaugeVec{
		"qps":         newGaugeVec("qps", "qps measurement", perfstatsLabels),
		"avg_latency": newGaugeVec("avg_latency", "avg_latency measurement", perfstatsLabels),
	}
)

func newExporterPerfstats() Exporter {
	return exporterPerfstats{
		perfstatsMetricsGauge: perfstatsGaugeVec,
	}
}

func (e exporterPerfstats) String() string {
	return "Exporter profstats"
}

func (e exporterPerfstats) Collect(ch chan<- prometheus.Metric) error {
	perfstatsData, _ := loadJson(config, "#perfstats")
	perfstats := PerfstatsResponse{}
	if err := json.Unmarshal(perfstatsData, &perfstats); err != nil {
		return err
	}

	for gaugeKey, gauge := range e.perfstatsMetricsGauge {
		for _, item := range perfstats.Items {
			name := item.Name

			if strings.Compare(gaugeKey, "qps") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.LastSecQPS))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.LastSecQPS))
			}

			if strings.Compare(gaugeKey, "avg_latency") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.TotalAvgLatencyUs))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.TotalAvgLatencyUs))
			}
		}

		gauge.Collect(ch)
	}

	return nil
}

func (e exporterPerfstats) Describe(ch chan<- *prometheus.Desc) {
	for _, perfstatsMetric := range e.perfstatsMetricsGauge {
		perfstatsMetric.Describe(ch)
	}
}
