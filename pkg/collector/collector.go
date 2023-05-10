package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"prom-logstash-exporter/constants"
	"prom-logstash-exporter/pkg/collector/node_stats"
	"prom-logstash-exporter/pkg/helpers"
	"prom-logstash-exporter/pkg/restclient"
	"sync"
)

type Collector struct {
	logstashClient   *LogstashClient
	metricsCollector *MetricsCollector
	mutex            sync.Mutex
}

func NewLogstashCollector(uri string) (*Collector, error) {
	client, err := NewLogstashClient(uri)
	if err != nil {
		return nil, err
	}

	metricsCollector := NewMetricsCollector()

	return &Collector{
		logstashClient:   client,
		metricsCollector: metricsCollector,
	}, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.metricsCollector.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // Protect metrics from concurrent collects
	defer c.mutex.Unlock()

	up := c.logstashClient.PerformScrape(c.metricsCollector, ch)
	c.metricsCollector.UpdateUp(up)
	c.metricsCollector.Collect(ch)
}

type LogstashClient struct {
	handler restclient.HTTPHandlerInterface
}

func NewLogstashClient(logstashURL string) (*LogstashClient, error) {
	parsedURL, err := helpers.ParseURI(logstashURL)
	if err != nil {
		return nil, err
	}

	handler := &restclient.HTTPHandler{
		Endpoint: fmt.Sprintf("%s/%s", parsedURL, constants.StatsPath),
	}

	return &LogstashClient{
		handler: handler,
	}, nil
}

func (c *LogstashClient) PerformScrape(mc *MetricsCollector, ch chan<- prometheus.Metric) (up float64) {
	mc.IncrementTotalScrapes()

	var stats node_stats.NodeStats
	err := restclient.GetMetrics(c.handler, &stats)
	if err != nil {
		logrus.WithError(err).Warnln("Can't scrape Logstash", constants.StatsPath)
		return 0
	}

	mc.UpdateLogstashStatus(stats)
	mc.UpdateLogstashInfo(stats, ch)

	mc.jvm.Collect(stats.JVM, ch)
	mc.event.Collect(stats.Event, ch)
	mc.process.Collect(stats.Process, ch)
	mc.pipelines.Collect(stats.Pipelines, ch)
	mc.pipelineConfig.Collect(stats.Pipeline, ch)
	mc.reloadsConfig.Collect(stats.Reloads, ch)

	return 1
}

type MetricsCollector struct {
	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter
	logstashStatus    prometheus.Gauge
	logstashInfo      *prometheus.Desc
	jvm               *node_stats.JVMCollector
	event             *node_stats.EventCollector
	process           *node_stats.ProcessCollector
	pipelines         *node_stats.PipelinesCollector
	pipelineConfig    *node_stats.PipelineConfigCollector
	reloadsConfig     *node_stats.ReloadsConfigCollector
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: constants.Namespace,
			Name:      "up",
			Help:      "Was the last scrape of logstash successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: constants.Namespace,
			Name:      "exporter_total_scrapes",
			Help:      "Current total logstash scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: constants.Namespace,
			Name:      "exporter_json_parse_failures",
			Help:      "Number of errors while parsing JSON.",
		}),
		logstashStatus: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: constants.Namespace,
			Name:      "status",
			Help:      "Logstash status: 0 for Green; 1 for Yellow; 2 for Red.",
		}),
		logstashInfo:   prometheus.NewDesc(prometheus.BuildFQName(constants.Namespace, "", "info"), "A metric with a constant '1' value labeled by version, http_address, name, id and ephemeral_id from Logstash instance.", []string{"version", "http_address", "name", "id", "ephemeral_id"}, nil),
		jvm:            node_stats.NewJVMCollector(),
		event:          node_stats.NewEventCollector(),
		process:        node_stats.NewProcessCollector(),
		pipelines:      node_stats.NewPipelinesCollector(),
		pipelineConfig: node_stats.NewPipelineConfigCollector(),
		reloadsConfig:  node_stats.NewReloadsConfigCollector(),
	}
}

func (mc *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(mc, ch)
}

func (mc *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- mc.up
	ch <- mc.totalScrapes
	ch <- mc.jsonParseFailures
	ch <- mc.logstashStatus
}

func (mc *MetricsCollector) UpdateUp(up float64) {
	mc.up.Set(up)
}

func (mc *MetricsCollector) IncrementTotalScrapes() {
	mc.totalScrapes.Inc()
}

func (mc *MetricsCollector) IncrementJsonParseFailures() {
	mc.jsonParseFailures.Inc()
}

func (mc *MetricsCollector) UpdateLogstashStatus(stats node_stats.NodeStats) {
	switch stats.Status {
	case "green":
		mc.logstashStatus.Set(0)
	case "yellow":
		mc.logstashStatus.Set(1)
	default:
		mc.logstashStatus.Set(2)
	}
}

func (mc *MetricsCollector) UpdateLogstashInfo(stats node_stats.NodeStats, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(mc.logstashInfo, prometheus.GaugeValue, 1.0, stats.Version, stats.HttpAddress, stats.Name, stats.ID, stats.EphemeralID)
}
