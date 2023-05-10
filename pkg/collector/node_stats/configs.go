package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
)

type PipelineConfigCollector struct {
	Workers    *prometheus.Desc
	BatchSize  *prometheus.Desc
	BatchDelay *prometheus.Desc
}

func NewPipelineConfigCollector() *PipelineConfigCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "pipeline_config")
	return &PipelineConfigCollector{
		Workers:    desc("workers", "The number of workers that will, in parallel, execute the filter and output stages of the pipeline."),
		BatchSize:  desc("batch_size", "The maximum number of events an individual worker thread will collect from inputs before attempting to execute its filters and outputs."),
		BatchDelay: desc("batch_delay_seconds", "How long to wait before dispatching an undersized batch to workers."),
	}
}

func (c *PipelineConfigCollector) Collect(p PipelineConfig, ch chan<- prometheus.Metric) {
	metrics := []struct {
		desc      *prometheus.Desc
		valueType prometheus.ValueType
		value     float64
		labels    []string
	}{
		{c.Workers, prometheus.GaugeValue, float64(p.Workers), nil},
		{c.BatchSize, prometheus.GaugeValue, float64(p.BatchSize), nil},
		{c.BatchDelay, prometheus.GaugeValue, float64(p.BatchDelay) / 1000.0, nil},
	}

	for _, m := range metrics {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}
}

type ReloadsConfigCollector struct {
	Failures  *prometheus.Desc
	Successes *prometheus.Desc
}

func NewReloadsConfigCollector() *ReloadsConfigCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "reloads_config")
	return &ReloadsConfigCollector{
		Failures:  desc("failures_total", "Number of failures during config reload."),
		Successes: desc("successes_total", "Number of successful config reloads."),
	}
}

func (c *ReloadsConfigCollector) Collect(p ReloadsConfig, ch chan<- prometheus.Metric) {
	metrics := []struct {
		desc      *prometheus.Desc
		valueType prometheus.ValueType
		value     float64
		labels    []string
	}{
		{c.Failures, prometheus.CounterValue, float64(p.Failures), nil},
		{c.Successes, prometheus.CounterValue, float64(p.Successes), nil},
	}

	for _, m := range metrics {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}
}
