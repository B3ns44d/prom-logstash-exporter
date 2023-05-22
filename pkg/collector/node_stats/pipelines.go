package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
	"strconv"
)

type PipelinesCollector struct {
	In                *prometheus.Desc
	Filtered          *prometheus.Desc
	Out               *prometheus.Desc
	Duration          *prometheus.Desc
	QueuePushDuration *prometheus.Desc

	InputConnections       *prometheus.Desc
	InputQueuePushDuration *prometheus.Desc
	InputOut               *prometheus.Desc

	FilterDuration *prometheus.Desc
	FilterIn       *prometheus.Desc
	FilterOut      *prometheus.Desc

	OutputDuration             *prometheus.Desc
	OutputIn                   *prometheus.Desc
	OutputOut                  *prometheus.Desc
	OutputSuccesses            *prometheus.Desc
	OutputNonRetryableFailures *prometheus.Desc

	EventsCount  *prometheus.Desc
	QueueSize    *prometheus.Desc
	MaxQueueSize *prometheus.Desc

	CapacityMaxUnreadEvents     *prometheus.Desc
	CapacityMaxQueueSizeInBytes *prometheus.Desc
	CapacityPageCapacityInBytes *prometheus.Desc
	CapacityQueueSizeInBytes    *prometheus.Desc

	//DeadLetterQueue
	DroppedEvents              *prometheus.Desc
	MaxQueueSizeInBytes        *prometheus.Desc
	DeadLetterQueueSizeInBytes *prometheus.Desc
}

func NewPipelinesCollector() *PipelinesCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "pipeline")
	return &PipelinesCollector{
		In:                desc("event_in_total", "The total number of events in.", "pipeline"),
		Filtered:          desc("event_filtered_total", "The total numbers of filtered.", "pipeline"),
		Out:               desc("event_out_total", "The total number of events out.", "pipeline"),
		Duration:          desc("event_duration_seconds_total", "The total process duration time in seconds.", "pipeline"),
		QueuePushDuration: desc("event_queue_push_duration_seconds_total", "The total in queue duration time in seconds.", "pipeline"),

		InputConnections:       desc("input_connections", "The current number of connections.", "pipeline", "id", "name"),
		InputQueuePushDuration: desc("input_queue_push_seconds_total", "The total in queue duration time in seconds", "pipeline", "id", "name"),
		InputOut:               desc("input_out_total", "The total number of events out.", "pipeline", "id", "name"),

		FilterDuration: desc("filter_duration_seconds_total", "The total process duration time in seconds", "pipeline", "id", "name", "index"),
		FilterIn:       desc("filter_in_total", "The total number of events in.", "pipeline", "id", "name", "index"),
		FilterOut:      desc("filter_out_total", "The total number of events out.", "pipeline", "id", "name", "index"),

		OutputDuration:             desc("output_duration_seconds_total", "The total process duration time in seconds", "pipeline", "id", "name"),
		OutputIn:                   desc("output_in_total", "The total number of events in.", "pipeline", "id", "name"),
		OutputOut:                  desc("output_out_total", "The total number of events out.", "pipeline", "id", "name"),
		OutputSuccesses:            desc("output_successes_total", "The total number of successful outputs.", "pipeline", "id", "name"),
		OutputNonRetryableFailures: desc("output_non_retryable_failures_total", "The total number of non-retryable output failures.", "pipeline", "id", "name"),

		EventsCount:  desc("queue_event_count", "The current events in queue.", "pipeline", "queue_type"),
		QueueSize:    desc("queue_size_bytes", "The current queue size in bytes.", "pipeline", "queue_type"),
		MaxQueueSize: desc("queue_max_size_bytes", "The max queue size in bytes.", "pipeline", "queue_type"),

		CapacityMaxUnreadEvents:     desc("capacity_max_unread_events", "The maximum number of unread events in capacity.", "pipeline", "queue_type"),
		CapacityMaxQueueSizeInBytes: desc("capacity_max_queue_size_bytes", "The maximum size of the capacity queue in bytes.", "pipeline", "queue_type"),
		CapacityPageCapacityInBytes: desc("page_capacity_bytes", "The capacity of a single page in bytes.", "pipeline", "queue_type"),
		CapacityQueueSizeInBytes:    desc("capacity_queue_size_bytes", "The current size of the queue capacity in bytes.", "pipeline", "queue_type"),

		DroppedEvents:              desc("dead_letter_queue_dropped_events_total", "The total number of dropped events in the dead letter queue.", "pipeline"),
		MaxQueueSizeInBytes:        desc("dead_letter_queue_max_queue_size_bytes", "The maximum size of the dead letter queue in bytes.", "pipeline"),
		DeadLetterQueueSizeInBytes: desc("dead_letter_queue_size_bytes", "The current size of the dead letter queue in bytes.", "pipeline"),
	}
}

type pipelineMetricData struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     float64
	labels    []string
}

func (c *PipelinesCollector) Collect(p map[string]Pipeline, ch chan<- prometheus.Metric) {
	for pipelineName, pipeline := range p {
		c.collectMetricsForPipeline(pipelineName, pipeline, ch)
	}
}

func (c *PipelinesCollector) collectMetricsForPipeline(pipelineName string, p Pipeline, ch chan<- prometheus.Metric) {
	eventMetrics := []pipelineMetricData{
		{c.In, prometheus.CounterValue, float64(p.Event.In), []string{pipelineName}},
		{c.Filtered, prometheus.CounterValue, float64(p.Event.Filtered), []string{pipelineName}},
		{c.Out, prometheus.CounterValue, float64(p.Event.Out), []string{pipelineName}},
		{c.Duration, prometheus.CounterValue, float64(p.Event.DurationInMillis) / 1000.0, []string{pipelineName}},
		{c.QueuePushDuration, prometheus.CounterValue, float64(p.Event.QueuePushDurationInMillis) / 1000.0, []string{pipelineName}},
	}

	queueMetrics := []pipelineMetricData{
		{c.EventsCount, prometheus.GaugeValue, float64(p.Queue.EventsCount), []string{pipelineName, p.Queue.Type}},
		{c.QueueSize, prometheus.CounterValue, float64(p.Queue.QueueSizeInBytes), []string{pipelineName, p.Queue.Type}},
		{c.MaxQueueSize, prometheus.CounterValue, float64(p.Queue.MaxQueueSizeInBytes), []string{pipelineName, p.Queue.Type}},
		{c.CapacityMaxUnreadEvents, prometheus.CounterValue, float64(p.Queue.Capacity.MaxUnreadEvents), []string{pipelineName, p.Queue.Type}},
		{c.CapacityMaxQueueSizeInBytes, prometheus.CounterValue, float64(p.Queue.Capacity.MaxQueueSizeInBytes), []string{pipelineName, p.Queue.Type}},
		{c.CapacityPageCapacityInBytes, prometheus.CounterValue, float64(p.Queue.Capacity.PageCapacityInBytes), []string{pipelineName, p.Queue.Type}},
		{c.CapacityQueueSizeInBytes, prometheus.CounterValue, float64(p.Queue.Capacity.QueueSizeInBytes), []string{pipelineName, p.Queue.Type}},
	}

	deadLetterQueueMetrics := []pipelineMetricData{
		{c.DroppedEvents, prometheus.CounterValue, float64(p.DeadLetterQueue.DroppedEvents), []string{pipelineName}},
		{c.MaxQueueSizeInBytes, prometheus.CounterValue, float64(p.DeadLetterQueue.MaxQueueSizeInBytes), []string{pipelineName}},
		{c.DeadLetterQueueSizeInBytes, prometheus.CounterValue, float64(p.DeadLetterQueue.QueueSizeInBytes), []string{pipelineName}},
	}
	var inputMetrics, filterMetrics, outputMetrics []pipelineMetricData

	for _, plugin := range p.Plugins.Inputs {
		inputMetrics = append(inputMetrics,
			pipelineMetricData{c.InputConnections, prometheus.GaugeValue, float64(plugin.CurrentConnections), []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.InputQueuePushDuration, prometheus.CounterValue, float64(plugin.Events.QueuePushDurationInMillis) / 1000.0, []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.InputOut, prometheus.CounterValue, float64(plugin.Events.Out), []string{pipelineName, plugin.ID, plugin.Name}},
		)
	}

	for idx, plugin := range p.Plugins.Filters {
		index := strconv.Itoa(idx)
		filterMetrics = append(filterMetrics,
			pipelineMetricData{c.FilterDuration, prometheus.CounterValue, float64(plugin.Events.DurationInMillis) / 1000.0, []string{pipelineName, plugin.ID, plugin.Name, index}},
			pipelineMetricData{c.FilterIn, prometheus.CounterValue, float64(plugin.Events.In), []string{pipelineName, plugin.ID, plugin.Name, index}},
			pipelineMetricData{c.FilterOut, prometheus.CounterValue, float64(plugin.Events.Out), []string{pipelineName, plugin.ID, plugin.Name, index}},
		)
	}

	for _, plugin := range p.Plugins.Outputs {
		outputMetrics = append(outputMetrics,
			pipelineMetricData{c.OutputDuration, prometheus.CounterValue, float64(plugin.Events.DurationInMillis) / 1000.0, []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.OutputIn, prometheus.CounterValue, float64(plugin.Events.In), []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.OutputOut, prometheus.CounterValue, float64(plugin.Events.Out), []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.OutputSuccesses, prometheus.CounterValue, float64(plugin.Documents.Successes), []string{pipelineName, plugin.ID, plugin.Name}},
			pipelineMetricData{c.OutputNonRetryableFailures, prometheus.CounterValue, float64(plugin.Documents.NonRetryableFailures), []string{pipelineName, plugin.ID, plugin.Name}},
		)
	}

	for _, m := range append(append(append(append(append(eventMetrics, queueMetrics...), deadLetterQueueMetrics...), inputMetrics...), filterMetrics...), outputMetrics...) {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}
}
