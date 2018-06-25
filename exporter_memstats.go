package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

func init() {
	RegisterExporter("memstats", newExporterMemstats)
}

type MemstatsTotal struct {
	DataSize    int `json:"data_size"`
	IndexesSize int `json:"indexes_size"`
	CacheSize   int `json:"cache_size"`
}

type MemstatsCacheDef struct {
	TotalSize     int `json:"total_size"`
	ItemsCount    int `json:"items_count"`
	EmptyCount    int `json:"empty_count"`
	HitCountLimit int `json:"hit_count_limit"`
}

type MemstatsJoinCache struct {
	MemstatsCacheDef
}

type MemstatsQueryCache struct {
	MemstatsCacheDef
}

type MemstatsIndexIdsetCache struct {
	MemstatsCacheDef
}

type MemstatsIndex struct {
	UniqKeysCount  int                     `json:"uniq_keys_count"`
	DataSize       int                     `json:"data_size,omitempty"`
	Name           string                  `json:"name"`
	IdsetPlainSize int                     `json:"idset_plain_size,omitempty"`
	IdsetCache     MemstatsIndexIdsetCache `json:"idset_cache,omitempty"`
	SortOrdersSize int                     `json:"sort_orders_size,omitempty"`
}

type MemstatsItems struct {
	Name            string             `json:"name"`
	ItemsCount      int                `json:"items_count"`
	DataSize        int                `json:"data_size"`
	UpdatedUnixNano int64              `json:"updated_unix_nano"`
	StorageOk       bool               `json:"storage_ok"`
	StoragePath     string             `json:"storage_path"`
	Total           MemstatsTotal      `json:"total"`
	JoinCache       MemstatsJoinCache  `json:"join_cache"`
	QueryCache      MemstatsQueryCache `json:"query_cache"`
	Indexes         []MemstatsIndex    `json:"indexes"`
	EmptyItemsCount int                `json:"empty_items_count,omitempty"`
}

type MemstatsResponse struct {
	Items      []MemstatsItems `json:"items"`
	TotalItems int             `json:"total_items"`
}

type exporterMemstats struct {
	memstatsMetricsGauge map[string]*prometheus.GaugeVec
}

var (
	memstatsLabels   = []string{"namespace"}
	memstatsGaugeVec = map[string]*prometheus.GaugeVec{
		"data_size":    newGaugeVec("data_size", "namespace data size", memstatsLabels),
		"indexes_size": newGaugeVec("indexes_size", "namespace indexes size", memstatsLabels),
		"caches_size":  newGaugeVec("caches_size", "namespace cache size", memstatsLabels),
	}
)

func newExporterMemstats() Exporter {
	return exporterMemstats{
		memstatsMetricsGauge: memstatsGaugeVec,
	}
}

func (e exporterMemstats) String() string {
	return "Exporter profstats"
}

func (e exporterMemstats) Collect(ch chan<- prometheus.Metric) error {
	memstatsData, _ := loadJson(config, "#memstats")
	memstats := MemstatsResponse{}
	if err := json.Unmarshal(memstatsData, &memstats); err != nil {
		return err
	}

	for gaugeKey, gauge := range e.memstatsMetricsGauge {
		for _, item := range memstats.Items {
			name := item.Name

			if strings.Compare(gaugeKey, "data_size") == 0 {
				gauge.WithLabelValues(name).Set(float64(item.Total.DataSize))
			}

			if strings.Compare(gaugeKey, "indexes_size") == 0 {
				gauge.WithLabelValues(name).Set(float64(item.Total.IndexesSize))
			}

			if strings.Compare(gaugeKey, "caches_size") == 0 {
				gauge.WithLabelValues(name).Set(float64(item.Total.CacheSize))
			}
		}

		gauge.Collect(ch)
	}

	return nil
}

func (e exporterMemstats) Describe(ch chan<- *prometheus.Desc) {
	for _, memstatsMetric := range e.memstatsMetricsGauge {
		memstatsMetric.Describe(ch)
	}
}
