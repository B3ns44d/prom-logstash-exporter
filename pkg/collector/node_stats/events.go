package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
)

type EventCollector struct {
	In                *prometheus.Desc
	Filtered          *prometheus.Desc
	Out               *prometheus.Desc
	Duration          *prometheus.Desc
	QueuePushDuration *prometheus.Desc
}

func NewEventCollector() *EventCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "event")
	return &EventCollector{
		In:                desc("in_total", "The total number of events in."),
		Filtered:          desc("filtered_total", "The total numbers of filtered."),
		Out:               desc("out_total", "The total number of events out."),
		Duration:          desc("duration_seconds_total", "The total process duration time in seconds."),
		QueuePushDuration: desc("queue_push_duration_seconds_total", "The total in queue duration time in seconds."),
	}
}

type eventMetricData struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     float64
	labels    []string
}

func (c *EventCollector) Collect(e Event, ch chan<- prometheus.Metric) {
	metrics := []eventMetricData{
		{c.In, prometheus.CounterValue, float64(e.In), nil},
		{c.Filtered, prometheus.CounterValue, float64(e.Filtered), nil},
		{c.Out, prometheus.CounterValue, float64(e.Out), nil},
		{c.Duration, prometheus.CounterValue, float64(e.DurationInMillis) / 1000.0, nil},
		{c.QueuePushDuration, prometheus.CounterValue, float64(e.QueuePushDurationInMillis) / 1000.0, nil},
	}

	for _, m := range metrics {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}
}
