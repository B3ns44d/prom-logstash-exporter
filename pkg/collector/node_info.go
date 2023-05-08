package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/restclient"
	"strconv"
)

type NodeInfoCollector struct {
	endpoint string

	NodeInfos *prometheus.Desc
	OsInfos   *prometheus.Desc
	JvmInfos  *prometheus.Desc
}

type metricData struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     float64
	labels    []string
}

func newInfoDesc(subsystem, name, help string, labelNames []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(constants.Namespace, subsystem, name),
		help,
		labelNames,
		nil,
	)
}

func sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labels ...string) {
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labels...)
	if err != nil {
		logrus.Errorf("Failed to create metric %s: %v", desc, err)
		return
	}
	ch <- metric
}

func NewNodeInfoCollector(logstashEndpoint string) (Collector, error) {
	const subsystem = "info"

	nodeInfos := newInfoDesc(subsystem, "node", "A metric with a constant '1' value labeled by Logstash version.", []string{"version"})
	osInfos := newInfoDesc(subsystem, "os", "A metric with a constant '1' value labeled by name, arch, version, and available_processors to the OS running Logstash.", []string{"name", "arch", "version", "available_processors"})
	jvmInfos := newInfoDesc(subsystem, "jvm", "A metric with a constant '1' value labeled by name, version, and vendor of the JVM running Logstash.", []string{"name", "version", "vendor"})

	return &NodeInfoCollector{
		endpoint:  logstashEndpoint,
		NodeInfos: nodeInfos,
		OsInfos:   osInfos,
		JvmInfos:  jvmInfos,
	}, nil
}

func (c *NodeInfoCollector) Collect(ch chan<- prometheus.Metric) error {
	stats, err := restclient.NodeInfo(c.endpoint)
	if err != nil {
		logrus.Error("Failed collecting info metrics: ", err)
		return err
	}

	metrics := []metricData{
		{c.NodeInfos, prometheus.GaugeValue, float64(1), []string{stats.Version}},
		{c.OsInfos, prometheus.GaugeValue, float64(1), []string{stats.Os.Name, stats.Os.Arch, stats.Os.Version, strconv.Itoa(stats.Os.AvailableProcessors)}},
		{c.JvmInfos, prometheus.GaugeValue, float64(1), []string{stats.Jvm.VmName, stats.Jvm.VmVersion, stats.Jvm.VmVendor}},
	}

	for _, m := range metrics {
		sendMetric(ch, m.desc, m.valueType, m.value, m.labels...)
	}

	return nil
}
