package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func init() {
	RegisterExporter("perfstats", newExporterPerfstats)
}

type PerfstatAction struct {
	TotalQueriesCount    int `json:"total_queries_count"`
	TotalAvgLatencyUs    int `json:"total_avg_latency_us"`
	TotalAvgLockTimeUs   int `json:"total_avg_lock_time_us"`
	LastSecQPS           int `json:"last_sec_qps"`
	LastSecAvgLockTimeUs int `json:"last_sec_avg_lock_time_us"`
	LastSecAvgLatencyUs  int `json:"last_sec_avg_latency_us"`
}

type PerfstatsResponse struct {
	Items []struct {
		Name    string         `json:"name"`
		Updates PerfstatAction `json:"updates"`
		Selects PerfstatAction `json:"selects"`
	} `json:"items"`
}

type exporterPerfstats struct {
	perfstatsMetricsGauge map[string]*prometheus.GaugeVec
}

var (
	perfstatsLabels   = []string{"namespace", "action"}
	perfstatsGaugeVec = map[string]*prometheus.GaugeVec{
		"total_queries_count":       newGaugeVec("total_queries_count", "total_queries_count measurement", perfstatsLabels),
		"total_avg_latency_us":      newGaugeVec("total_avg_latency_us", "total_avg_latency_us measurement", perfstatsLabels),
		"total_avg_lock_time_us":    newGaugeVec("total_avg_lock_time_us", "total_avg_lock_time_us measurement", perfstatsLabels),
		"last_sec_qps":              newGaugeVec("last_sec_qps", "last_sec_qps measurement", perfstatsLabels),
		"last_sec_avg_lock_time_us": newGaugeVec("last_sec_avg_lock_time_us", "last_sec_avg_lock_time_us measurement", perfstatsLabels),
		"last_sec_avg_latency_us":   newGaugeVec("last_sec_avg_latency_us", "last_sec_avg_latency_us measurement", perfstatsLabels),
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
	perfstatsData, _ := apiGetQuery(config, "#perfstats")
	perfstats := PerfstatsResponse{}
	if err := json.Unmarshal(perfstatsData, &perfstats); err != nil {
		return err
	}

	for gaugeKey, gauge := range e.perfstatsMetricsGauge {
		for _, item := range perfstats.Items {
			name := item.Name

			if strings.Compare(gaugeKey, "total_queries_count") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.TotalQueriesCount))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.TotalQueriesCount))
			}

			if strings.Compare(gaugeKey, "total_avg_latency_us") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.TotalAvgLatencyUs))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.TotalAvgLatencyUs))
			}

			if strings.Compare(gaugeKey, "total_avg_lock_time_us") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.TotalAvgLockTimeUs))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.TotalAvgLockTimeUs))
			}

			if strings.Compare(gaugeKey, "last_sec_qps") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.LastSecQPS))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.LastSecQPS))
			}

			if strings.Compare(gaugeKey, "last_sec_avg_lock_time_us") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.LastSecAvgLockTimeUs))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.LastSecAvgLockTimeUs))
			}

			if strings.Compare(gaugeKey, "last_sec_avg_latency_us") == 0 {
				gauge.WithLabelValues(name, "select").Set(float64(item.Selects.LastSecAvgLatencyUs))
				gauge.WithLabelValues(name, "update").Set(float64(item.Updates.LastSecAvgLatencyUs))
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
