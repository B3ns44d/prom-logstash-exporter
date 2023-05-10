package node_stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/helpers"
)

type JVMCollector struct {
	threadsCount         *prometheus.Desc
	heapUsedRatio        *prometheus.Desc
	heapCommittedInBytes *prometheus.Desc
	heapUsedInBytes      *prometheus.Desc
	poolUsedBytes        *prometheus.Desc
	poolCommittedBytes   *prometheus.Desc
	poolMaxBytes         *prometheus.Desc
	gc                   *prometheus.Desc
}

func NewJVMCollector() *JVMCollector {
	desc := helpers.NewDescFQ(constants.Namespace, "jvm")
	return &JVMCollector{
		threadsCount:         desc("threads_count", "Current JVM thread count."),
		heapUsedRatio:        desc("heap_used_ratio", "Current JVM heap usage ratio."),
		heapCommittedInBytes: desc("heap_committed_bytes", "Current JVM heap committed size"),
		heapUsedInBytes:      desc("heap_used_bytes", "Current JVM heap used size"),
		poolUsedBytes:        desc("memory_pool_used_bytes", "Current JVM heap pool used size", "pool"),
		poolCommittedBytes:   desc("memory_pool_committed_bytes", "Current JVM heap pool committed size", "pool"),
		poolMaxBytes:         desc("memory_pool_max_bytes", "Current JVM heap pool max size", "pool"),
		gc:                   desc("gc_collection_duration_seconds", "GC collection duration.", "collector"),
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
		{c.threadsCount, prometheus.GaugeValue, float64(jvm.Threads.Count), nil},
		{c.heapUsedRatio, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedPercent) / 100.0, nil},
		{c.heapCommittedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapCommittedInBytes), nil},
		{c.heapUsedInBytes, prometheus.GaugeValue, float64(jvm.Mem.HeapUsedInBytes), nil},
	}

	poolLabels := []string{"young", "survivor", "old"}
	poolMetrics := []jvmMetricData{
		{c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.UsedInBytes), []string{poolLabels[0]}},
		{c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.UsedInBytes), []string{poolLabels[1]}},
		{c.poolUsedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.UsedInBytes), []string{poolLabels[2]}},
		{c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.CommittedInBytes), []string{poolLabels[0]}},
		{c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.CommittedInBytes), []string{poolLabels[1]}},
		{c.poolCommittedBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.CommittedInBytes), []string{poolLabels[2]}},
		{c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Young.MaxInBytes), []string{poolLabels[0]}},
		{c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Survivor.MaxInBytes), []string{poolLabels[1]}},
		{c.poolMaxBytes, prometheus.GaugeValue, float64(jvm.Mem.Pools.Old.MaxInBytes), []string{poolLabels[2]}},
	}

	for _, m := range append(metrics, poolMetrics...) {
		ch <- prometheus.MustNewConstMetric(m.desc, m.valueType, m.value, m.labels...)
	}

	ch <- prometheus.MustNewConstSummary(c.gc, jvm.GC.Collectors.Young.CollectionCount, float64(jvm.GC.Collectors.Young.CollectionTimeInMillis)/1000.0, nil, "young")
	ch <- prometheus.MustNewConstSummary(c.gc, jvm.GC.Collectors.Old.CollectionCount, float64(jvm.GC.Collectors.Old.CollectionTimeInMillis)/1000.0, nil, "old")
}
