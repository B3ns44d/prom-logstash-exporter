package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
)

type JVMCollector struct {
	ThreadsCount         *prometheus.Desc
	HeapUsedRatio        *prometheus.Desc
	HeapCommittedInBytes *prometheus.Desc
	HeapUsedInBytes      *prometheus.Desc
	PoolUsedBytes        *prometheus.Desc
	PoolCommittedBytes   *prometheus.Desc
	PoolMaxBytes         *prometheus.Desc
	GC                   *prometheus.Desc
}

func NewJVMCollector() *JVMCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "jvm")
	return &JVMCollector{
		ThreadsCount:         desc("threads_count", "Current JVM thread count."),
		HeapUsedRatio:        desc("heap_used_ratio", "Current JVM heap usage ratio."),
		HeapCommittedInBytes: desc("heap_committed_bytes", "Current JVM heap committed size"),
		HeapUsedInBytes:      desc("heap_used_bytes", "Current JVM heap used size"),
		PoolUsedBytes:        desc("memory_pool_used_bytes", "Current JVM heap pool used size", "pool"),
		PoolCommittedBytes:   desc("memory_pool_committed_bytes", "Current JVM heap pool committed size", "pool"),
		PoolMaxBytes:         desc("memory_pool_max_bytes", "Current JVM heap pool max size", "pool"),
		GC:                   desc("gc_collection_duration_seconds", "GC collection duration.", "collector"),
	}
}

type jvmMetricData struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     float64
	labels    []string
}

func (c *JVMCollector) Collect(jvm JVM, ch chan<- prometheus.Metric) {

	metrics := []jvmMetricData{
		{c.ThreadsCount, prometheus.GaugeValue, float64(jvm.Threads.Count), nil},
		{c.HeapUsedRatio, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedPercent) / 100.0, nil},
		{c.HeapCommittedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapCommittedInBytes), nil},
		{c.HeapUsedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedInBytes), nil},
	}

	poolLabels := []string{"young", "survivor", "old"}
	poolMetrics := []jvmMetricData{
		{c.PoolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.UsedInBytes), []string{poolLabels[0]}},
		{c.PoolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.UsedInBytes), []string{poolLabels[1]}},
		{c.PoolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.UsedInBytes), []string{poolLabels[2]}},
		{c.PoolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.CommittedInBytes), []string{poolLabels[0]}},
		{c.PoolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.CommittedInBytes), []string{poolLabels[1]}},
		{c.PoolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.CommittedInBytes), []string{poolLabels[2]}},
		{c.PoolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.MaxInBytes), []string{poolLabels[0]}},
		{c.PoolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.MaxInBytes), []string{poolLabels[1]}},
		{c.PoolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.MaxInBytes), []string{poolLabels[2]}},
	}

	for _, m := range append(metrics, poolMetrics...) {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}

	ch <- prometheus.MustNewConstSummary(c.GC, jvm.GC.Collectors.Young.CollectionCount, float64(jvm.GC.Collectors.Young.CollectionTimeInMillis)/1000.0, nil, "young")
	ch <- prometheus.MustNewConstSummary(c.GC, jvm.GC.Collectors.Old.CollectionCount, float64(jvm.GC.Collectors.Old.CollectionTimeInMillis)/1000.0, nil, "old")
}
