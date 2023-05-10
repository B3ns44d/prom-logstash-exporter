package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
)

type ProcessCollector struct {
	OpenFileDescriptors *prometheus.Desc
	MaxFileDescriptors  *prometheus.Desc
	TotalVirtualMemory  *prometheus.Desc
	ProcessTime         *prometheus.Desc
	CPUUsage            *prometheus.Desc
	LoadAverage         *prometheus.Desc
}

func NewProcessCollector() *ProcessCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "process")
	return &ProcessCollector{
		OpenFileDescriptors: desc("open_file_descriptors", "Current open file descriptors"),
		MaxFileDescriptors:  desc("max_file_descriptors", "Max file descriptors"),
		TotalVirtualMemory:  desc("total_virtual_memory_bytes", "Was the used virtual memory."),
		ProcessTime:         desc("process_time_seconds", "Was the total process time."),
		CPUUsage:            desc("cpu_usage_ratio", "Was the CPU usage"),
		LoadAverage:         desc("load_average", "Was the system load average", "load"),
	}
}

type processMetricData struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     float64
	labels    []string
}

func (c *ProcessCollector) Collect(p Process, ch chan<- prometheus.Metric) {
	metrics := []processMetricData{
		{c.OpenFileDescriptors, prometheus.GaugeValue, float64(p.OpenFileDescriptors), nil},
		{c.MaxFileDescriptors, prometheus.GaugeValue, float64(p.MaxFileDescriptors), nil},
		{c.TotalVirtualMemory, prometheus.GaugeValue, float64(p.Mem.TotalVirtualInBytes), nil},
		{c.ProcessTime, prometheus.CounterValue, float64(p.CPU.TotalInMillis) / 1000.0, nil},
		{c.CPUUsage, prometheus.GaugeValue, float64(p.CPU.Percent) / 100.0, nil},
	}

	loadLabels := []string{"1", "5", "15"}
	loadMetrics := []processMetricData{
		{c.LoadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load1, []string{loadLabels[0]}},
		{c.LoadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load5, []string{loadLabels[1]}},
		{c.LoadAverage, prometheus.GaugeValue, p.CPU.LoadAverage.Load15, []string{loadLabels[2]}},
	}

	for _, m := range append(metrics, loadMetrics...) {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}
}
