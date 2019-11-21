package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector generic collector type
type Collector struct {
	firstMetric *prometheus.Desc
}

// NewCollector returns the collector
func NewCollector() *Collector {
	return &Collector{
		firstMetric: prometheus.NewDesc("my_first_metric_name",
			"Description of the first metric",
			nil, nil,
		),
	}
}

// Describe to satisfy the collector interface.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.firstMetric
}

// Collect metrics from ...
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	collectedValue := 42.0
	ch <- prometheus.MustNewConstMetric(c.firstMetric, prometheus.CounterValue, collectedValue)
}
