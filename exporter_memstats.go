package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	RegisterExporter("memstats", newExporterMemstats)
}

type MemStatsTotal struct {
	DataSize    int `json:"data_size"`
	IndexesSize int `json:"indexes_size"`
	CacheSize   int `json:"cache_size"`
}

type MemStatsCacheDef struct {
	TotalSize     int `json:"total_size"`
	ItemsCount    int `json:"items_count"`
	EmptyCount    int `json:"empty_count"`
	HitCountLimit int `json:"hit_count_limit"`
}

type MemStatsJoinCache struct {
	MemStatsCacheDef
}

type MemStatsQueryCache struct {
	MemStatsCacheDef
}

type MemStatsIndexIdsetCache struct {
	MemStatsCacheDef
}

type MemStatsIndex struct {
	UniqKeysCount  int                     `json:"uniq_keys_count"`
	DataSize       int                     `json:"data_size,omitempty"`
	Name           string                  `json:"name"`
	IdsetPlainSize int                     `json:"idset_plain_size,omitempty"`
	IdsetCache     MemStatsIndexIdsetCache `json:"idset_cache,omitempty"`
	SortOrdersSize int                     `json:"sort_orders_size,omitempty"`
}

type MemStatsItems struct {
	Name            string             `json:"name"`
	ItemsCount      int                `json:"items_count"`
	DataSize        int                `json:"data_size"`
	UpdatedUnixNano int64              `json:"updated_unix_nano"`
	StorageOk       bool               `json:"storage_ok"`
	StoragePath     string             `json:"storage_path"`
	Total           MemStatsTotal      `json:"total"`
	JoinCache       MemStatsJoinCache  `json:"join_cache"`
	QueryCache      MemStatsQueryCache `json:"query_cache"`
	Indexes         []MemStatsIndex    `json:"indexes"`
	EmptyItemsCount int                `json:"empty_items_count,omitempty"`
}

type ReindexerMemStatsResponse struct {
	Items      []MemStatsItems `json:"items"`
	TotalItems int             `json:"total_items"`
}

type exporterMemStats struct {
	memstatsMetricsGauge []*prometheus.GaugeVec
}

var (
	memstatsGaugeVec = []*prometheus.GaugeVec{
		newGaugeVec("memstats_data_size", "TODO", nil),
		newGaugeVec("memstats_updated_unix_nano", "TODO", nil),
		newGaugeVec("memstats_storage_ok", "TODO", nil),
		newGaugeVec("memstats_storage_path", "TODO", nil),

		newGaugeVec("memstats_total_data_size", "TODO", nil),
		newGaugeVec("memstats_total_indexes_size", "TODO", nil),
		newGaugeVec("memstats_total_cache_size", "TODO", nil),

		newGaugeVec("memstats_join_cache_total_size", "TODO", nil),
		newGaugeVec("memstats_join_cache_items_count", "TODO", nil),
		newGaugeVec("memstats_join_cache_empty_count", "TODO", nil),
		newGaugeVec("memstats_join_cache_hit_count_limit", "TODO", nil),

		newGaugeVec("memstats_query_cache_total_size", "TODO", nil),
		newGaugeVec("memstats_query_cache_items_count", "TODO", nil),
		newGaugeVec("memstats_query_cache_empty_count", "TODO", nil),
		newGaugeVec("memstats_query_cache_hit_count_limit", "TODO", nil),
		newGaugeVec("memstats_empty_items_count", "TODO", nil),
	}
)

func newExporterMemstats() Exporter {
	return exporterMemStats{
		memstatsMetricsGauge: memstatsGaugeVec,
	}
}

func (e exporterMemStats) String() string {
	return "Exporter memstats"
}

func (e exporterMemStats) Collect(ch chan<- prometheus.Metric) error {
	memstatsData, _ := loadJson(config, "#memstats")
	memStats := ReindexerMemStatsResponse{}
	if err := json.Unmarshal(memstatsData, &memStats); err != nil {
		return err
	}

	for _, item := range memStats.Items {
		// name := item.`Name

		for _, gauge := range e.memstatsMetricsGauge {
			gauge.With(labels).Set(float64(item.ItemsCount))
			gauge.Collect(ch)
		}

	}

	return nil
}

func (e exporterMemStats) Describe(ch chan<- *prometheus.Desc) {
	for _, memStatsMetric := range e.memstatsMetricsGauge {
		memStatsMetric.Describe(ch)
	}
}
